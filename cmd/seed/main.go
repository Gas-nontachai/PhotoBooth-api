package main

import (
	"context"
	"log"

	"go-ddd-clean/internal/infrastructure/config"
	"go-ddd-clean/internal/infrastructure/db"
	"go-ddd-clean/internal/infrastructure/seeder"
)

func main() {
	cfg := config.LoadConfig()
	database := db.ConnectDB(cfg.DB_DSN)

	if err := seeder.Run(context.Background(), database); err != nil {
		log.Fatal("❌ Failed to seed database:", err)
	}
	log.Println("✅ Seed data applied successfully")
}
