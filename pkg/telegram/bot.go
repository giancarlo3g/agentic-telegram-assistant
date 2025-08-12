package telegram

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot handles all Telegram bot interactions
type Bot struct {
	bot *tgbotapi.BotAPI
}

// NewBot creates a new Telegram bot instance
func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %v", err)
	}

	log.Printf("Telegram bot created successfully: %s", bot.Self.UserName)
	return &Bot{bot: bot}, nil
}

// SendMessage sends a message to a specific chat
func (t *Bot) SendMessage(chatID int64, text string) error {
	log.Printf("Sending message to chat %d: %s", chatID, text)

	// Check if message is too long for Telegram (max 4096 characters)
	if len(text) > 4096 {
		log.Printf("Message too long (%d chars), splitting into multiple messages", len(text))
		return t.sendLongMessage(chatID, text)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

// sendLongMessage splits a long message into multiple parts and sends them
func (t *Bot) sendLongMessage(chatID int64, text string) error {
	const maxLength = 4000 // Leave some buffer for safety

	// Split the message into chunks
	var messages []string
	start := 0

	for start < len(text) {
		end := start + maxLength
		if end > len(text) {
			end = len(text)
		}

		// Try to break at a newline if possible
		if end < len(text) {
			lastNewline := strings.LastIndex(text[start:end], "\n")
			if lastNewline > 0 && lastNewline > maxLength/2 {
				end = start + lastNewline + 1
			}
		}

		messages = append(messages, text[start:end])
		start = end
	}

	// Send each chunk
	for i, message := range messages {
		msg := tgbotapi.NewMessage(chatID, message)
		if i < len(messages)-1 {
			msg.Text += "\n\n[Message continued...]"
		}

		_, err := t.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message part %d/%d: %v", i+1, len(messages), err)
		}

		// Small delay between messages to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("Successfully sent long message in %d parts", len(messages))
	return nil
}

// SendMessageWithKeyboard sends a message with an inline keyboard
func (t *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message with keyboard: %v", err)
	}
	return nil
}

// GetUpdatesChan returns the updates channel for the bot
func (t *Bot) GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return t.bot.GetUpdatesChan(u)
}

// GetBotInfo returns information about the bot
func (t *Bot) GetBotInfo() tgbotapi.User {
	return t.bot.Self
}

// AnswerCallbackQuery answers a callback query
func (t *Bot) AnswerCallbackQuery(callbackQueryID string, text string) error {
	callback := tgbotapi.NewCallback(callbackQueryID, text)
	_, err := t.bot.Request(callback)
	if err != nil {
		return fmt.Errorf("failed to answer callback query: %v", err)
	}
	return nil
}

// EditMessageText edits an existing message
func (t *Bot) EditMessageText(chatID int64, messageID int, newText string) error {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, newText)
	_, err := t.bot.Send(edit)
	if err != nil {
		return fmt.Errorf("failed to edit message: %v", err)
	}
	return nil
}

// DeleteMessage deletes a message
func (t *Bot) DeleteMessage(chatID int64, messageID int) error {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := t.bot.Request(deleteMsg)
	if err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}

// CreateInlineKeyboard creates an inline keyboard with the given buttons
func CreateInlineKeyboard(buttons [][]tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// CreateInlineKeyboardButton creates a single inline keyboard button
func CreateInlineKeyboardButton(text, callbackData string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(text, callbackData)
}

// CreateCalendarKeyboard creates a keyboard for calendar navigation
func CreateCalendarKeyboard() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			CreateInlineKeyboardButton("Today", "calendar_today"),
			CreateInlineKeyboardButton("Tomorrow", "calendar_tomorrow"),
		},
	}
	return CreateInlineKeyboard(buttons)
}
