package state

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/joyboy1210/tex/internal/models"
	"github.com/joyboy1210/tex/internal/twilio"
	"github.com/joyboy1210/tex/internal/utils"
)

const (
	StateStart          = "START"
	StateBrowsing       = "BROWSING"
	StateViewingProduct = "VIEWING_PRODUCT"
	StateAwaitingQty    = "AWAITING_QTY"
	StateCartDecision   = "CART_DECISION"
)

func ProcessMessage(ctx context.Context, phone, message string) {
	cleanInput := strings.TrimSpace(strings.ToLower(message))
	// _ = cleanInput
	currentState, exists := GetState(phone)
	if !exists {
		user, err := models.GetUser(phone)
		if err != nil {
			log.Printf("new user with phone number: %s", phone)
			_, err := models.CreateUser(phone, StateStart)
			if err != nil {
				log.Printf("ERROR: failed to create user for %s : %v", phone, err)
				return
			}
			currentState = StateStart
		} else {
			currentState = user.CurrentState
		}
		SetState(phone, currentState)
	}

	switch currentState {
	case StateStart:
		if cleanInput == "1" {
			TransitionState(phone, StateBrowsing)
			handleBrowsing(ctx, phone)
		} else if cleanInput == "2" {
			orders, err := models.GetActiveOrders(phone)

			if err != nil || len(orders) == 0 {
				twilio.SendMessage(ctx, phone, "You don't have any active orders right now! \n\nType '1' to start shopping.")
				return
			}

			var msg strings.Builder
			msg.WriteString("  *Your Active Orders*\n\n")

			for _, order := range orders {
				msg.WriteString(fmt.Sprintf("   *Order #TEX-%d*\n", order.ID))
				msg.WriteString(fmt.Sprintf("   Total: $%.2f\n", order.TotalPrice))
				msg.WriteString(fmt.Sprintf("   Status: *%s*\n\n", order.Status))
			}

			msg.WriteString("Type '0' or 'menu' to return to the main menu.")
			twilio.SendMessage(ctx, phone, msg.String())
		} else if cleanInput == "3" {

			var msg strings.Builder
			msg.WriteString("  *How to use Tex Bot*\n\n")
			msg.WriteString("• Reply *1* to browse our premium football jerseys.\n")
			msg.WriteString("• Reply *2* to check the status of your active orders.\n")
			msg.WriteString("• Type *0* at almost any time to go back to the previous menu.\n")
			msg.WriteString("Ready? Type *1* to start shopping!")

			twilio.SendMessage(ctx, phone, msg.String())

		} else {

			handleStart(ctx, phone)
		}
	case StateBrowsing:
		if cleanInput == "0" {
			TransitionState(phone, StateStart)
			handleStart(ctx, phone)
			return
		}
		choice, err := strconv.Atoi(cleanInput)
		if err != nil {
			twilio.SendMessage(ctx, phone, "Please reply with a valid number from the menu.")
			handleBrowsing(ctx, phone)
			return
		}
		categories := GetCategories()
		index := choice - 1
		if index < 0 || index >= len(categories) {
			twilio.SendMessage(ctx, phone, "Please reply with a valid number from the menu.")
			handleBrowsing(ctx, phone)
			return
		}
		selectedCategory := categories[index]
		utils.SetSessionCategory(phone, selectedCategory.ID)
		TransitionState(phone, StateViewingProduct)

		sendProductMenu(ctx, phone, selectedCategory.ID)

	case StateViewingProduct:
		if cleanInput == "0" {
			TransitionState(phone, StateBrowsing)
			handleBrowsing(ctx, phone)
			return
		}
		session, exists := utils.GetSession(phone)
		if !exists {
			TransitionState(phone, StateStart)
			handleStart(ctx, phone)
			return
		}
		choice, err := strconv.Atoi(cleanInput)
		if err != nil {
			twilio.SendMessage(ctx, phone, "Please reply with a valid number from the menu.")
			sendProductMenu(ctx, phone, session.CategoryId)
			return
		}
		products := GetProductsByCategoryID(session.CategoryId)
		index := choice - 1
		if index < 0 || index >= len(products) {
			twilio.SendMessage(ctx, phone, "Please reply with a valid number from the menu.")
			sendProductMenu(ctx, phone, session.CategoryId)
			return
		}
		selectedProduct := products[index]
		utils.SetSessionProduct(phone, selectedProduct.ID)
		TransitionState(phone, StateAwaitingQty)
		twilio.SendMessage(ctx, phone, fmt.Sprintf("You selected *%s*.\nPlease reply with the quantity you want to order.\n *Reply with 0 to go back*", selectedProduct.Name))

	case StateAwaitingQty:
		if cleanInput == "0" {
			TransitionState(phone, StateBrowsing)
			handleBrowsing(ctx, phone)
			return
		}
		session, exits := utils.GetSession(phone)
		if !exits {
			TransitionState(phone, StateStart)
			handleStart(ctx, phone)
			return
		}
		qty, err := strconv.Atoi(cleanInput)
		if err != nil || qty <= 0 {
			twilio.SendMessage(ctx, phone, "Please reply with a valid quantity (a positive number).")
			return
		}
		err = models.AddToCart(phone, session.ProductId, qty)
		if err != nil {
			log.Printf("ERROR: failed to add product %d with quantity %d to cart for user %s : %v", session.ProductId, qty, phone, err)
			twilio.SendMessage(ctx, phone, "Sorry, there was an error adding the product to your cart. Please try again.")
			return
		}
		TransitionState(phone, StateCartDecision)

		var message strings.Builder
		message.WriteString(fmt.Sprintf("Added %d item(s) to your cart!\n\n", qty))
		message.WriteString("*What would you like to do next?*\n")
		message.WriteString("1. Checkout \n")
		message.WriteString("2. Keep Shopping \n")
		message.WriteString("3. View Cart \n")
		message.WriteString("4. Clear Cart ")

		twilio.SendMessage(ctx, phone, message.String())

	case StateCartDecision:
		switch cleanInput {
		case "1":
			cart, err := models.GetCart(phone)
			if err != nil || len(cart) == 0 {
				twilio.SendMessage(ctx, phone, "Your cart is already empty. Please add items before checking out.")
				TransitionState(phone, StateBrowsing)
				handleBrowsing(ctx, phone)
				return
			}

			var total float64
			var receipt strings.Builder
			receipt.WriteString("*Order Receipt*\n\n")

			for _, item := range cart {
				product, exists := GetProductByID(item.ProductID)
				if exists {
					lineTotal := product.Price * float64(item.Quantity)
					total += lineTotal
					receipt.WriteString(fmt.Sprintf("%dx %s - $%.2f\n", item.Quantity, product.Name, lineTotal))
				}
			}

			order, err := models.CreateOrder(phone, cart, total)
			if err != nil {
				log.Printf("ERROR: failed to create order for user %s : %v", phone, err)
				twilio.SendMessage(ctx, phone, "Sorry, there was an error processing your order. Please try again.")
				return
			}

			receipt.WriteString(fmt.Sprintf("\n*Total: $%.2f*\n\n", total))
			receipt.WriteString(fmt.Sprintf("Your order has been placed! Your Order ID is #%d.\n\nType 'hi' to start a new session.", order.ID))
			twilio.SendMessage(ctx, phone, receipt.String())
			models.ClearCart(phone)
			utils.ClearSession(phone)
			TransitionState(phone, StateStart)

		case "2":
			TransitionState(phone, StateBrowsing)
			handleBrowsing(ctx, phone)

		case "3":
			cart, err := models.GetCart(phone)
			if err != nil || len(cart) == 0 {
				twilio.SendMessage(ctx, phone, "Your cart is currently empty. Please add some products before viewing your cart.")
				return
			}

			var total float64
			var sb strings.Builder
			sb.WriteString("  *Your Current Cart*\n\n")

			for _, item := range cart {
				product, exists := GetProductByID(item.ProductID)
				if exists {
					lineTotal := product.Price * float64(item.Quantity)
					total += lineTotal
					sb.WriteString(fmt.Sprintf("%dx %s ($%.2f each)\n", item.Quantity, product.Name, product.Price))
				}
			}

			sb.WriteString(fmt.Sprintf("\n*Total: $%.2f*\n\n", total))
			sb.WriteString("Reply *1* to Checkout, *2* to Keep Shopping, or *4* to Clear Cart.")

			twilio.SendMessage(ctx, phone, sb.String())

		case "4":
			err := models.ClearCart(phone)
			if err != nil {
				log.Printf("ERROR: failed to clear cart for user %s : %v", phone, err)
				twilio.SendMessage(ctx, phone, "Sorry, there was an error clearing your cart. Please try again.")
				return
			}
			twilio.SendMessage(ctx, phone, "Your cart has been cleared. You can start adding products again!")
			TransitionState(phone, StateBrowsing)
			handleBrowsing(ctx, phone)

		default:
			twilio.SendMessage(ctx, phone, "Please reply with a valid option from the menu.")
		}

	default:
		log.Printf("Unhandled state: %s for user %s", currentState, phone)
		TransitionState(phone, StateStart)
		handleStart(ctx, phone)
	}
}

func handleStart(ctx context.Context, phone string) {
	log.Printf("Welcome message sent to %s", phone)
	log.Printf("routing to main menu")

	message := "Welcome to our store! Please choose an option:\n1. Browse Products\n2. Track Order\n3. Help"

	err := twilio.SendMessage(ctx, phone, message)
	if err != nil {
		log.Printf("ERROR: failed to send welcome message to %s : %v", phone, err)
		return
	}
	log.Printf("welcome message sent to %s", phone)

}

func handleBrowsing(ctx context.Context, phone string) {
	log.Printf("routing %s to category menu", phone)

	categories := GetCategories()

	if len(categories) == 0 {
		log.Printf("No categories available to display to %s", phone)
		twilio.SendMessage(ctx, phone, "Sorry, no categories are available at the moment.")
		return
	}

	var message strings.Builder

	message.WriteString("*Store Catalog*\n\nReply with a number to explore:\n\n")
	for i, cat := range categories {
		message.WriteString(fmt.Sprintf("%d. %s\n", i+1, cat.Name))

	}
	message.WriteString("\n 0. Back to Main Menu")
	err := twilio.SendMessage(ctx, phone, message.String())
	if err != nil {
		log.Printf("ERROR: failed to send category menu to %s : %v", phone, err)
		return
	}
	log.Printf("category menu sent to %s", phone)
}

func sendProductMenu(ctx context.Context, phone string, categoryID uint) {
	log.Printf("routing %s to product menu for category %d", phone, categoryID)

	products := GetProductsByCategoryID(categoryID)
	if len(products) == 0 {
		twilio.SendMessage(ctx, phone, "This category is currently empty! Check back later. 🚧")
		TransitionState(phone, StateBrowsing)
		handleBrowsing(ctx, phone)
		return
	}

	var message strings.Builder

	message.WriteString("*Products:*\n\nReply with a number to select:\n\n")

	for i, prod := range products {
		message.WriteString(fmt.Sprintf("%d. *%s*\n", i+1, prod.Name))
		message.WriteString(fmt.Sprintf("   %s\n", prod.Description))
		message.WriteString(fmt.Sprintf("    $%.2f\n\n", prod.Price))
	}
	message.WriteString("\n 0. Back to Categories")

	err := twilio.SendMessage(ctx, phone, message.String())
	if err != nil {
		log.Printf("ERROR: failed to send product menu to %s : %v", phone, err)
		return
	}
	log.Printf("product menu sent to %s for category %d", phone, categoryID)
}
