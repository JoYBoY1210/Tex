package api

import (
	"net/http"

	"github.com/joyboy1210/tex/internal/api/handlers"
)


func RegisterRoutes(mux *http.ServeMux){
	mux.HandleFunc("POST /webhook", handlers.WebhookHandler)
}