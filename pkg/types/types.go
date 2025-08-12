package types

import "time"

// AIResponse represents the AI agent's response
type AIResponse struct {
	Action     string `json:"action"`
	Message    string `json:"message"`
	EventID    string `json:"event_id,omitempty"`
	EventTitle string `json:"event_title,omitempty"`
	EventDate  string `json:"event_date,omitempty"`
	EventTime  string `json:"event_time,omitempty"`
	EventDesc  string `json:"event_description,omitempty"`
	EventLoc   string `json:"event_location,omitempty"`
	// For complex requests, AI can specify multiple actions
	Actions []AIAction `json:"actions,omitempty"`
}

// AIAction represents a single action the AI wants to perform
type AIAction struct {
	Action     string `json:"action"`
	EventDate  string `json:"event_date,omitempty"`
	EventTitle string `json:"event_title,omitempty"`
	EventTime  string `json:"event_time,omitempty"`
	EventDesc  string `json:"event_description,omitempty"`
	EventLoc   string `json:"event_location,omitempty"`
}

// CalendarEvent represents a calendar event
type CalendarEvent struct {
	ID          string    `json:"id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location"`
}

// Interaction represents a single interaction with the AI
type Interaction struct {
	UserID      int64     `json:"user_id"`
	Timestamp   time.Time `json:"timestamp"`
	UserMessage string    `json:"user_message"`
	AIResponse  string    `json:"ai_response"`
	Action      string    `json:"action,omitempty"`
}
