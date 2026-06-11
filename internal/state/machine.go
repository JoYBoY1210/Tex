package state

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/joyboy1210/tex/internal/models"
	"github.com/joyboy1210/tex/internal/twilio"
)

const (
	StateStart          = "START"
	StateBrowsing       = "BROWSING"
	StateViewingProduct = "VIEWING_PRODUCT"
	StateAwaitingQty    = "AWAITING_QTY"
	StateCartDecision   = "CART_DECISION"
)

var (
	sessionMut   sync.RWMutex
	userSessions = make(map[string]uint)
)

func setSessionCategory(phone string, categoryId uint) {
	sessionMut.Lock()
	defer sessionMut.Unlock()
	userSessions[phone] = categoryId
}

func getSessionCategory(phone string) (uint, bool) {
	sessionMut.RLock()
	defer sessionMut.RUnlock()
	catId, exists := userSessions[phone]
	return catId, exists
}

func TransitionState(phone, newState string) error {
	err := models.UpdateUserState(phone, newState)
	if err != nil {
		log.Printf("ERROR: database state update failed for %s : %v", phone, err)
		return err
	}
	SetState(phone, newState)
	log.Printf("state Transition: %s is now in %s", phone, newState)
	return nil
}

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
			twilio.SendMessage(ctx, phone, "Tracking system nahi banaya lol")
		} else {
			handleStart(ctx, phone)
		}
	case StateBrowsing:
		if cleanInput == "0" {
			TransitionState(phone, StateStart)
			handleStart(ctx, phone)
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
		setSessionCategory(phone, selectedCategory.ID)
		TransitionState(phone, StateViewingProduct)

		sendProductMenu(ctx, phone, selectedCategory.ID)

	case StateViewingProduct:
		log.Printf("user is viewing product")
		twilio.SendMessage(ctx, phone, "Jersey selected\n\n(Cart logic coming right after this)")

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
