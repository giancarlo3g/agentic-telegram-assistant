package main

import (
	"context"
	"fmt"
	"log"

	calapi "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"calendar-assistant-bot/pkg/ai"
	calendarpkg "calendar-assistant-bot/pkg/calendar"
	"calendar-assistant-bot/pkg/config"
	"calendar-assistant-bot/pkg/database"
	"calendar-assistant-bot/pkg/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents the main bot that coordinates all components
type Bot struct {
	aiAgent     *ai.Agent
	telegramBot *telegram.Bot
	config      *config.Config
}

// NewBot creates a new bot instance with all components
func NewBot(cfg *config.Config) (*Bot, error) {
	log.Printf("Creating bot with config: Telegram=%s, OpenAI=%s, GoogleCreds=%s, CalendarID=%s",
		config.MaskToken(cfg.TelegramToken), config.MaskToken(cfg.OpenAIKey), cfg.GoogleCreds, cfg.CalendarID)

	// Create OpenAI service
	openaiService := ai.NewOpenAIService(cfg.OpenAIKey)
	log.Printf("OpenAI service created successfully")

	// Create Google Calendar service
	ctx := context.Background()
	log.Printf("Creating Google Calendar service with credentials file: %s", cfg.GoogleCreds)
	calendarService, err := calapi.NewService(ctx, option.WithCredentialsFile(cfg.GoogleCreds))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %v", err)
	}
	log.Printf("Google Calendar service created successfully")

	// Create Google Calendar tool
	calendarTool := calendarpkg.NewService(calendarService, cfg.CalendarID)
	log.Printf("Google Calendar tool created successfully")

	// Create Telegram bot
	telegramBot, err := telegram.NewBot(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %v", err)
	}
	log.Printf("Telegram bot created successfully: %s", telegramBot.GetBotInfo().UserName)

	// Create database
	database, err := database.NewDatabase("./data")
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}
	log.Printf("Database created successfully")

	// Create AI agent
	aiAgent := ai.NewAgent(openaiService, calendarTool, telegramBot, database)
	log.Printf("AI agent created successfully")

	return &Bot{
		aiAgent:     aiAgent,
		telegramBot: telegramBot,
		config:      cfg,
	}, nil
}

// handleMessage processes incoming Telegram messages
func (b *Bot) handleMessage(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	message := update.Message.Text
	chatID := update.Message.Chat.ID

	log.Printf("Received message from user %d (chatID %d): %s", userID, chatID, message)

	// Process message through AI agent
	if err := b.aiAgent.ProcessUserMessage(userID, chatID, message); err != nil {
		log.Printf("Error processing message for user %d: %v", userID, err)
	}
}

// startBot starts the Telegram bot
func (b *Bot) startBot() error {
	updates := b.telegramBot.GetUpdatesChan()

	log.Printf("Bot started. Listening for messages...")

	for update := range updates {
		go b.handleMessage(update)
	}

	return nil
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	bot, err := NewBot(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Starting calendar assistant bot...")

	if err := bot.startBot(); err != nil {
		log.Fatalf("Bot error: %v", err)
	}
}
