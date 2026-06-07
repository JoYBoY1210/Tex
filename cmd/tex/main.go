package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joyboy1210/tex/internal/api"
	"github.com/joyboy1210/tex/internal/db"
)

func main() {
	Db, err := db.InitDB("tex.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	db.Migrate(Db)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	api.InitServer(ctx)
}
