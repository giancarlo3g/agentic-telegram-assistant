package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"calendar-assistant-bot/pkg/types"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIService handles all interactions with OpenAI
type OpenAIService struct {
	client *openai.Client
}

// NewOpenAIService creates a new OpenAI service instance
func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(apiKey),
	}
}

// ProcessMessage sends a user message to OpenAI and returns the AI response
func (o *OpenAIService) ProcessMessage(userContext, message string) (*types.AIResponse, error) {
	systemPrompt := `You are a calendar assistant. Your responsibilities include creating, getting, and deleting events in the user's calendar.

Available actions:
- getEvents: Get events for a specific date
- delEvents: Delete a specific event (requires event ID)
- makeEvent: Create a new event
- updtEvent: Update an event (requires event ID)

Current date/time: ` + time.Now().Format("2006-01-02 15:04:05") + `

You are intelligent and flexible. You can:
- Interpret natural language date requests ("last week", "this weekend", "next month")
- Convert relative dates to actual dates (YYYY-MM-DD format)
- Handle complex requests by making multiple API calls if needed
- Use "today", "tomorrow", "yesterday" as keywords

IMPORTANT: When a user asks to "get events for today" or similar, you MUST respond with action="getEvents" and event_date="today". Do NOT respond with action="None".
You can provide current date and time if asked but make sure it includes time zone which is UTC.

For complex requests like "what did I do last week?", you can either:
1. Make a single getEvents call with the calculated date range, OR
2. Use the actions array to make multiple getEvents calls for different days

If no duration is specified for an event, assume it will be one hour.

Respond with a JSON object containing:
- action: one of the available actions (for simple requests). If no actions are needed, action should be "None".
- message: response to user
- event_id: if deleting/updating (get this from getEvents first)
- event_title, event_date, event_time, event_description, event_location: if creating/updating
- actions: array of actions for complex requests (optional)

Example responses:
{"action": "getEvents", "message": "I'll get events for today", "event_date": "today"}
{"action": "getEvents", "message": "I'll get events for last week (Aug 5-11)", "event_date": "2025-08-05"}
{"actions": [{"action": "getEvents", "event_date": "2025-08-05"}, {"action": "getEvents", "event_date": "2026-08-06"}], "message": "I'll get events for Monday and Tuesday of last week"}`

	userPrompt := message
	if userContext != "" {
		userPrompt = userContext + "\n\nUser: " + message
	}

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	// Try to parse JSON response
	var aiResp types.AIResponse
	if err := json.Unmarshal([]byte(content), &aiResp); err != nil {
		// If not JSON, treat as simple message
		aiResp = types.AIResponse{
			Action:  "message",
			Message: content,
		}
	}

	// Validate required fields for specific actions
	if aiResp.Action == "getEvents" && aiResp.EventDate == "" {
		// If getEvents action but no date specified, default to today
		aiResp.EventDate = "today"
		log.Printf("AI response missing event_date for getEvents action, defaulting to 'today'")
	}

	return &aiResp, nil
}
