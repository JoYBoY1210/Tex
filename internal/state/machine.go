package state

import (
	"log"
	"strings"

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

func ProcessMessage(phone, message string) {
	cleanInput := strings.TrimSpace(strings.ToLower(message))
	_ = cleanInput
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
		handleStart(phone)
	case StateBrowsing:
		log.Printf("user is browsing")
	default:
		log.Printf("Unhandled state: %s for user %s", currentState, phone)
		TransitionState(phone, StateStart)
		handleStart(phone)
	}
}

func handleStart(phone string) {
	log.Printf("Welcome message sent to %s", phone)
	log.Printf("routing to main menu")
}
