package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/joyboy1210/tex/internal/state"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	defer w.WriteHeader(http.StatusOK)

	if err := r.ParseForm(); err != nil {
		fmt.Println("failed to parse form: ", err)
		return
	}
	sender := r.FormValue("From")
	messageBody := strings.TrimSpace(r.FormValue("Body"))

	phoneNumber := strings.Replace(sender, "whatsapp:", "", 1)

	if phoneNumber == "" {
		log.Println("received webhook with empty phone number")
		return
	}

	log.Println("received msg from %s: %s \n", phoneNumber, messageBody)

	ctx := r.Context()

	state.ProcessMessage(ctx, phoneNumber, messageBody)
}
