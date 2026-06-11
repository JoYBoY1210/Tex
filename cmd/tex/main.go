package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/joyboy1210/tex/internal/api"
	"github.com/joyboy1210/tex/internal/db"
	"github.com/joyboy1210/tex/internal/models"
)

func main() {
	Db, err := db.InitDB("tex.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	err = godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env")
	}
	models.InitDB(Db)
	db.Migrate(Db)
	db.SeedDb(Db)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	api.InitServer(ctx)
}
