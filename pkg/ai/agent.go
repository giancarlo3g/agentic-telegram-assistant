package ai

import (
	"calendar-assistant-bot/pkg/calendar"
	"calendar-assistant-bot/pkg/database"
	"calendar-assistant-bot/pkg/telegram"
	"calendar-assistant-bot/pkg/types"
	"fmt"
	"log"
)

// Agent coordinates between all tools and handles the main logic
type Agent struct {
	openaiService   *OpenAIService
	calendarService *calendar.Service
	telegramBot     *telegram.Bot
	database        *database.Database
}

// NewAgent creates a new AI agent instance
func NewAgent(openaiService *OpenAIService, calendarService *calendar.Service, telegramBot *telegram.Bot, database *database.Database) *Agent {
	return &Agent{
		openaiService:   openaiService,
		calendarService: calendarService,
		telegramBot:     telegramBot,
		database:        database,
	}
}

// ProcessUserMessage handles a complete user message flow
func (a *Agent) ProcessUserMessage(userID int64, chatID int64, message string) error {
	log.Printf("Processing message from user %d: %s", userID, message)

	// Get user context from database
	userContext := a.database.GetUserContext(userID, 10)

	// Send message to OpenAI for processing
	aiResponse, err := a.openaiService.ProcessMessage(userContext, message)
	if err != nil {
		log.Printf("AI processing error for user %d: %v", userID, err)
		errorMsg := "Sorry, I encountered an error processing your request. Please try again."
		if err := a.telegramBot.SendMessage(chatID, errorMsg); err != nil {
			log.Printf("Failed to send error message: %v", err)
		}
		return err
	}

	log.Printf("AI response for user %d: Action=%s, Message=%s, EventDate='%s'", userID, aiResponse.Action, aiResponse.Message, aiResponse.EventDate)

	// Store the interaction in database
	log.Printf("About to save interaction to database for user %d", userID)
	if err := a.database.AddInteraction(userID, message, aiResponse.Message, aiResponse.Action); err != nil {
		log.Printf("Failed to store interaction for user %d: %v", userID, err)
	}
	log.Printf("Successfully saved interaction to database for user %d", userID)

	log.Printf("About to execute AI action for user %d", userID)

	// Execute the AI's decision
	log.Printf("Calling executeAIAction for user %d", userID)
	response, err := a.executeAIAction(userID, aiResponse)
	log.Printf("executeAIAction returned for user %d: response='%s', err=%v", userID, response, err)
	if err != nil {
		log.Printf("Error executing AI action for user %d: %v", userID, err)
		response = fmt.Sprintf("Error executing action: %v", err)
	}

	// Send response to user
	log.Printf("About to send response to Telegram for user %d: %s", userID, response)
	if err := a.telegramBot.SendMessage(chatID, response); err != nil {
		log.Printf("Failed to send response to user %d: %v", userID, err)
		return err
	}
	log.Printf("Successfully sent response to Telegram for user %d", userID)

	return nil
}

// executeAIAction executes the action decided by the AI
func (a *Agent) executeAIAction(userID int64, aiResponse *types.AIResponse) (string, error) {
	log.Printf("executeAIAction ENTRY for user %d", userID)
	var response string
	log.Printf("Executing action for user %d, Action=%s, Message=%s, EventDate=%s", userID, aiResponse.Action, aiResponse.Message, aiResponse.EventDate)

	// Check if AI wants to perform multiple actions
	if len(aiResponse.Actions) > 0 {
		log.Printf("AI requested %d actions for user %d", len(aiResponse.Actions), userID)
		response = aiResponse.Message + "\n\n"

		for i, action := range aiResponse.Actions {
			log.Printf("Executing action %d/%d: %s", i+1, len(aiResponse.Actions), action.Action)

			switch action.Action {
			case "getEvents":
				events, err := a.calendarService.GetEvents(action.EventDate)
				if err != nil {
					log.Printf("Error getting events for user %d: %v", userID, err)
					response += fmt.Sprintf("Error getting events for %s: %v\n", action.EventDate, err)
				} else if len(events) == 0 {
					response += fmt.Sprintf("No events found for %s.\n", action.EventDate)
				} else {
					// Limit the number of events shown to prevent extremely long messages
					maxEvents := 15 // Slightly lower for multiple actions
					eventsToShow := events
					if len(events) > maxEvents {
						eventsToShow = events[:maxEvents]
						response += fmt.Sprintf("Events for %s (showing first %d of %d):\n", action.EventDate, maxEvents, len(events))
					} else {
						response += fmt.Sprintf("Events for %s:\n", action.EventDate)
					}

					for _, event := range eventsToShow {
						response += fmt.Sprintf("• %s (%s - %s)",
							event.Summary,
							event.Start.Format("15:04"),
							event.End.Format("15:04"))
						if event.Location != "" {
							response += fmt.Sprintf(" - %s", event.Location)
						}
						response += "\n"
					}

					if len(events) > maxEvents {
						response += fmt.Sprintf("... and %d more events.\n", len(events)-maxEvents)
					}
					response += "\n"
				}
			default:
				response += fmt.Sprintf("Unknown action: %s\n", action.Action)
			}
		}
	} else {
		// Single action (existing logic)
		switch aiResponse.Action {
		case "getEvents":
			log.Printf("Getting events for user %d, date: '%s' (length: %d)", userID, aiResponse.EventDate, len(aiResponse.EventDate))
			log.Printf("About to call Google Calendar API for user %d", userID)
			events, err := a.calendarService.GetEvents(aiResponse.EventDate)
			log.Printf("Google Calendar API call completed for user %d, err=%v, events count=%d", userID, err, len(events))
			if err != nil {
				log.Printf("Error getting events for user %d: %v", userID, err)
				response = fmt.Sprintf("Error getting events: %v", err)
			} else if len(events) == 0 {
				log.Printf("No events found for user %d on %s", userID, aiResponse.EventDate)
				response = fmt.Sprintf("No events found for %s.", aiResponse.EventDate)
			} else {
				log.Printf("Found %d events for user %d on %s", len(events), userID, aiResponse.EventDate)

				// Limit the number of events shown to prevent extremely long messages
				maxEvents := 20
				eventsToShow := events
				if len(events) > maxEvents {
					eventsToShow = events[:maxEvents]
					response = fmt.Sprintf("Events for %s (showing first %d of %d):\n", aiResponse.EventDate, maxEvents, len(events))
				} else {
					response = fmt.Sprintf("Events for %s:\n", aiResponse.EventDate)
				}

				for _, event := range eventsToShow {
					response += fmt.Sprintf("• %s (%s - %s)",
						event.Summary,
						event.Start.Format("15:04"),
						event.End.Format("15:04"))
					if event.Location != "" {
						response += fmt.Sprintf(" - %s", event.Location)
					}
					response += "\n"
				}

				if len(events) > maxEvents {
					response += fmt.Sprintf("\n... and %d more events. Use a more specific date range to see fewer events.", len(events)-maxEvents)
				}
			}

		case "makeEvent":
			log.Printf("Creating event for user %d: %s on %s at %s", userID, aiResponse.EventTitle, aiResponse.EventDate, aiResponse.EventTime)
			err := a.calendarService.CreateEvent(aiResponse.EventTitle, aiResponse.EventDate, aiResponse.EventTime, aiResponse.EventDesc, aiResponse.EventLoc)
			if err != nil {
				log.Printf("Error creating event for user %d: %v", userID, err)
				response = fmt.Sprintf("Error creating event: %v", err)
			} else {
				log.Printf("Successfully created event for user %d", userID)
				response = fmt.Sprintf("Event '%s' created successfully for %s at %s",
					aiResponse.EventTitle, aiResponse.EventDate, aiResponse.EventTime)
			}

		case "delEvents":
			if aiResponse.EventID == "" {
				log.Printf("User %d tried to delete event without specifying ID", userID)
				response = "Please specify an event ID to delete. Use 'getEvents' first to see available events."
			} else {
				log.Printf("Deleting event %s for user %d", aiResponse.EventID, userID)
				err := a.calendarService.DeleteEvent(aiResponse.EventID)
				if err != nil {
					log.Printf("Error deleting event %s for user %d: %v", aiResponse.EventID, userID, err)
					response = fmt.Sprintf("Error deleting event: %v", err)
				} else {
					log.Printf("Successfully deleted event %s for user %d", aiResponse.EventID, userID)
					response = "Event deleted successfully."
				}
			}

		case "updtEvent":
			if aiResponse.EventID == "" {
				log.Printf("User %d tried to update event without specifying ID", userID)
				response = "Please specify an event ID to update. Use 'getEvents' first to see available events."
			} else {
				log.Printf("Updating event %s for user %d", aiResponse.EventID, userID)
				err := a.calendarService.UpdateEvent(aiResponse.EventID, aiResponse.EventTitle, aiResponse.EventDate, aiResponse.EventTime, aiResponse.EventDesc, aiResponse.EventLoc)
				if err != nil {
					log.Printf("Error updating event %s for user %d: %v", aiResponse.EventID, userID, err)
					response = fmt.Sprintf("Error updating event: %v", err)
				} else {
					log.Printf("Successfully updated event %s for user %d", aiResponse.EventID, userID)
					response = "Event updated successfully."
				}
			}

		case "None":
			log.Printf("AI requested 'None' action for user %d", userID)
			response = aiResponse.Message

		case "message":
			log.Printf("AI sent a simple message for user %d", userID)
			response = aiResponse.Message

		default:
			log.Printf("No specific action for user %d, using AI message: %s", userID, aiResponse.Message)
			response = aiResponse.Message
		}
	}

	// Add logging to see what response we're about to return
	log.Printf("Final response for user %d: %s", userID, response)

	return response, nil
}

// HandleCalendarCallback handles calendar navigation callbacks
func (a *Agent) HandleCalendarCallback(userID int64, chatID int64, callbackData string) error {
	log.Printf("Handling calendar callback for user %d: %s", userID, callbackData)

	var response string

	switch callbackData {
	case "calendar_today":
		response = "Here are your events for today:"
		events, err := a.calendarService.GetEvents("today")
		if err != nil {
			response = fmt.Sprintf("Error getting events: %v", err)
		} else if len(events) == 0 {
			response = "No events found for today."
		} else {
			// Limit events to prevent long messages
			maxEvents := 15
			eventsToShow := events
			if len(events) > maxEvents {
				eventsToShow = events[:maxEvents]
				response += fmt.Sprintf("\n(showing first %d of %d):\n", maxEvents, len(events))
			} else {
				response += "\n"
			}

			for _, event := range eventsToShow {
				response += fmt.Sprintf("• %s (%s - %s)",
					event.Summary,
					event.Start.Format("15:04"),
					event.End.Format("15:04"))
				if event.Location != "" {
					response += fmt.Sprintf(" - %s", event.Location)
				}
				response += "\n"
			}

			if len(events) > maxEvents {
				response += fmt.Sprintf("... and %d more events.", len(events)-maxEvents)
			}
		}

	case "calendar_tomorrow":
		response = "Here are your events for tomorrow:"
		events, err := a.calendarService.GetEvents("tomorrow")
		if err != nil {
			response = fmt.Sprintf("Error getting events: %v", err)
		} else if len(events) == 0 {
			response = "No events found for tomorrow."
		} else {
			// Limit events to prevent long messages
			maxEvents := 15
			eventsToShow := events
			if len(events) > maxEvents {
				eventsToShow = events[:maxEvents]
				response += fmt.Sprintf("\n(showing first %d of %d):\n", maxEvents, len(events))
			} else {
				response += "\n"
			}

			for _, event := range eventsToShow {
				response += fmt.Sprintf("• %s (%s - %s)",
					event.Summary,
					event.Start.Format("15:04"),
					event.End.Format("15:04"))
				if event.Location != "" {
					response += fmt.Sprintf(" - %s", event.Location)
				}
				response += "\n"
			}

			if len(events) > maxEvents {
				response += fmt.Sprintf("... and %d more events.", len(events)-maxEvents)
			}
		}

	default:
		response = "Calendar navigation not implemented yet."
	}

	// Send response to user
	if err := a.telegramBot.SendMessage(chatID, response); err != nil {
		log.Printf("Failed to send calendar callback response: %v", err)
		return err
	}

	return nil
}

// GetUserStats retrieves user interaction statistics
func (a *Agent) GetUserStats(userID int64) map[string]interface{} {
	return a.database.GetUserStats(userID)
}

// CleanupOldInteractions removes old interactions from the database
func (a *Agent) CleanupOldInteractions(daysOld int) error {
	return a.database.Cleanup(daysOld)
}
