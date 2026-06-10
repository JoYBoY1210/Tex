package state

import (
	"context"
	"log"
	"strings"

	"github.com/joyboy1210/tex/internal/api/twilio"
	"github.com/joyboy1210/tex/internal/models"
)

const (
	StateStart          = "START"
	StateBrowsing       = "BROWSING"
	StateViewingProduct = "VIEWING_PRODUCT"
	StateAwaitingQty    = "AWAITING_QTY"
	StateCartDecision   = "CART_DECISION"
)

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
		log.Printf("User %s is actively browsing. Input: %s", phone, cleanInput)
		handleBrowsing(ctx, phone)
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
	message := "Here are our products:\n1. Product A\n2. Product B\n3. Product C\nPlease reply with the product number to view details."
	err := twilio.SendMessage(ctx, phone, message)
	if err != nil {
		log.Printf("ERROR: failed to send browsing message to %s : %v", phone, err)
		return
	}
	log.Printf("browsing message sent to %s", phone)
}
