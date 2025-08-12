package calendar

import (
	"context"
	"fmt"
	"log"
	"time"

	"calendar-assistant-bot/pkg/types"

	"google.golang.org/api/calendar/v3"
)

// Service handles all Google Calendar interactions
type Service struct {
	service    *calendar.Service
	calendarID string
}

// NewService creates a new Google Calendar service instance
func NewService(service *calendar.Service, calendarID string) *Service {
	tool := &Service{
		service:    service,
		calendarID: calendarID,
	}

	// Test the connection
	log.Printf("Testing Google Calendar connection...")
	events, err := tool.GetEvents("today")
	if err != nil {
		log.Printf("Warning: Google Calendar connection test failed: %v", err)
	} else {
		log.Printf("Google Calendar connection test successful, found %d events for today", len(events))
	}

	return tool
}

// GetEvents retrieves events from Google Calendar for a specific date
func (s *Service) GetEvents(dateStr string) ([]types.CalendarEvent, error) {
	// Parse date and set time range
	var startTime, endTime time.Time
	var err error

	// Handle empty date string by defaulting to today
	if dateStr == "" {
		dateStr = "today"
	}

	if dateStr == "today" {
		startTime = time.Now().Truncate(24 * time.Hour)
		endTime = startTime.Add(24 * time.Hour)
	} else if dateStr == "tomorrow" {
		startTime = time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)
		endTime = startTime.Add(24 * time.Hour)
	} else if dateStr == "yesterday" {
		startTime = time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour)
		endTime = startTime.Add(24 * time.Hour)
	} else {
		// Try to parse specific date
		startTime, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %v", err)
		}
		endTime = startTime.Add(24 * time.Hour)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	events, err := s.service.Events.List(s.calendarID).
		Context(ctx).
		TimeMin(startTime.Format(time.RFC3339)).
		TimeMax(endTime.Format(time.RFC3339)).
		OrderBy("startTime").
		SingleEvents(true).
		Do()

	if err != nil {
		return nil, fmt.Errorf("failed to get events: %v", err)
	}

	var calendarEvents []types.CalendarEvent
	for _, event := range events.Items {
		start, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		end, _ := time.Parse(time.RFC3339, event.End.DateTime)

		calendarEvents = append(calendarEvents, types.CalendarEvent{
			ID:          event.Id,
			Summary:     event.Summary,
			Description: event.Description,
			Start:       start,
			End:         end,
			Location:    event.Location,
		})
	}

	return calendarEvents, nil
}

// GetEventsInRange retrieves events from Google Calendar within a date range
func (s *Service) GetEventsInRange(startDate, endDate string) ([]types.CalendarEvent, error) {
	// Parse start and end dates
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	// Add one day to end date to include the full end date
	endTime = endTime.Add(24 * time.Hour)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	events, err := s.service.Events.List(s.calendarID).
		Context(ctx).
		TimeMin(startTime.Format(time.RFC3339)).
		TimeMax(endTime.Format(time.RFC3339)).
		OrderBy("startTime").
		SingleEvents(true).
		Do()

	if err != nil {
		return nil, fmt.Errorf("failed to get events: %v", err)
	}

	var calendarEvents []types.CalendarEvent
	for _, event := range events.Items {
		start, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		end, _ := time.Parse(time.RFC3339, event.End.DateTime)

		calendarEvents = append(calendarEvents, types.CalendarEvent{
			ID:          event.Id,
			Summary:     event.Summary,
			Description: event.Description,
			Start:       start,
			End:         end,
			Location:    event.Location,
		})
	}

	return calendarEvents, nil
}

// CreateEvent creates a new calendar event
func (s *Service) CreateEvent(title, dateStr, timeStr, description, location string) error {
	// Parse date and time
	dateTimeStr := dateStr + " " + timeStr
	startTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		return fmt.Errorf("invalid date/time format: %v", err)
	}

	endTime := startTime.Add(1 * time.Hour) // Default 1 hour duration

	event := &calendar.Event{
		Summary:     title,
		Description: description,
		Location:    location,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = s.service.Events.Insert(s.calendarID, event).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create event: %v", err)
	}

	return nil
}

// UpdateEvent updates an existing calendar event
func (s *Service) UpdateEvent(eventID, title, dateStr, timeStr, description, location string) error {
	// Parse date and time
	dateTimeStr := dateStr + " " + timeStr
	startTime, err := time.Parse("2006-01-02 15:04", dateTimeStr)
	if err != nil {
		return fmt.Errorf("invalid date/time format: %v", err)
	}

	endTime := startTime.Add(1 * time.Hour) // Default 1 hour duration

	event := &calendar.Event{
		Summary:     title,
		Description: description,
		Location:    location,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = s.service.Events.Update(s.calendarID, eventID, event).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to update event: %v", err)
	}

	return nil
}

// DeleteEvent deletes a calendar event
func (s *Service) DeleteEvent(eventID string) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.service.Events.Delete(s.calendarID, eventID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}

	return nil
}
