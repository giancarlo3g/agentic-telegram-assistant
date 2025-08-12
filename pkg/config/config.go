package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	TelegramToken string
	OpenAIKey     string
	GoogleCreds   string
	CalendarID    string
	Port          string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	log.Printf("Loading configuration...")

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		OpenAIKey:     os.Getenv("OPENAI_API_KEY"),
		GoogleCreds:   os.Getenv("GOOGLE_CREDENTIALS_FILE"),
		CalendarID:    os.Getenv("GOOGLE_CALENDAR_ID"),
		Port:          os.Getenv("PORT"),
	}

	if config.Port == "" {
		config.Port = "8080"
		log.Printf("Using default port: %s", config.Port)
	}

	log.Printf("Configuration loaded:")
	log.Printf("  Telegram Token: %s", MaskToken(config.TelegramToken))
	log.Printf("  OpenAI Key: %s", MaskToken(config.OpenAIKey))
	log.Printf("  Google Credentials: %s", config.GoogleCreds)
	log.Printf("  Calendar ID: %s", config.CalendarID)
	log.Printf("  Port: %s", config.Port)

	// Validate required config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	log.Printf("Configuration validation passed")
	return config, nil
}

// Validate checks if all required configuration values are present
func (c *Config) Validate() error {
	if c.TelegramToken == "" {
		return fmt.Errorf("TELEGRAM_TOKEN is required")
	}
	if c.OpenAIKey == "" {
		return fmt.Errorf("OPENAI_API_KEY is required")
	}
	if c.GoogleCreds == "" {
		return fmt.Errorf("GOOGLE_CREDENTIALS_FILE is required")
	}
	if c.CalendarID == "" {
		return fmt.Errorf("GOOGLE_CALENDAR_ID is required")
	}
	return nil
}

// MaskToken masks sensitive tokens for logging
func MaskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
