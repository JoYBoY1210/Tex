package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func InitServer(ctx context.Context) {
	mux := http.NewServeMux()

	server := http.Server{
		Addr:              ":6123",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		fmt.Printf("Listening on %s\n", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	<-ctx.Done()
	log.Println("Shutting down server")
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutDownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}
