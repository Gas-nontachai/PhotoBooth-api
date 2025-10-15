package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	DB_DSN           string
	BoothTokenSecret string
}

func LoadConfig() *Config {
	godotenv.Load()

	cfg := &Config{
		AppPort:          os.Getenv("APP_PORT"),
		DB_DSN:           os.Getenv("DB_DSN"),
		BoothTokenSecret: os.Getenv("BOOTH_TOKEN_SECRET"),
	}

	if cfg.AppPort == "" || cfg.DB_DSN == "" || cfg.BoothTokenSecret == "" {
		log.Fatal("Missing required environment variables")
	}

	return cfg
}
