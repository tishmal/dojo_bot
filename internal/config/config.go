package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Tgtoken  string
	MongoURI string
	ChanURL  string
	AppURL   string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Warning: .env file not loaded: %v", err)
	}

	env := map[string]string{
		"TELEGRAM_TOKEN": os.Getenv("TELEGRAM_TOKEN"),
		"MONGODB_URI":    os.Getenv("MONGODB_URI"),
		"CHANNEL_URL":    os.Getenv("CHANNEL_URL"),
		"APP_URL":        os.Getenv("APP_URL"),
	}

	for key, value := range env {
		if value == "" {
			return nil, fmt.Errorf("missing required env variable: %s", key)
		}
	}

	cfg := &Config{
		Tgtoken:  os.Getenv("TELEGRAM_TOKEN"),
		MongoURI: os.Getenv("MONGODB_URI"),
		ChanURL:  os.Getenv("CHANNEL_URL"),
		AppURL:   os.Getenv("APP_URL"),
	}

	return cfg, nil
}
