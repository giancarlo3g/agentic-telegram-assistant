package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"calendar-assistant-bot/pkg/types"
)

// Database handles storage and retrieval of AI interactions
type Database struct {
	filePath     string
	mutex        sync.RWMutex
	interactions map[int64][]types.Interaction
}

// NewDatabase creates a new database instance
func NewDatabase(dataDir string) (*Database, error) {
	db := &Database{
		filePath:     filepath.Join(dataDir, "interactions.json"),
		interactions: make(map[int64][]types.Interaction),
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	// Load existing interactions
	if err := db.loadInteractions(); err != nil {
		log.Printf("Warning: Could not load existing interactions: %v", err)
	}

	return db, nil
}

// AddInteraction stores a new interaction
func (d *Database) AddInteraction(userID int64, userMessage, aiResponse, action string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	interaction := types.Interaction{
		UserID:      userID,
		Timestamp:   time.Now(),
		UserMessage: userMessage,
		AIResponse:  aiResponse,
		Action:      action,
	}

	// Add to in-memory storage
	if d.interactions[userID] == nil {
		d.interactions[userID] = make([]types.Interaction, 0)
	}
	d.interactions[userID] = append(d.interactions[userID], interaction)

	// Keep only last 50 interactions per user
	if len(d.interactions[userID]) > 50 {
		d.interactions[userID] = d.interactions[userID][len(d.interactions[userID])-50:]
	}

	// Persist to disk
	log.Printf("Saving interactions for user %d", userID)
	log.Printf("Interactions: %v", d.interactions[userID])
	return d.saveInteractions()
}

// GetUserInteractions retrieves interactions for a specific user
func (d *Database) GetUserInteractions(userID int64, limit int) []types.Interaction {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	interactions := d.interactions[userID]
	if len(interactions) == 0 {
		return []types.Interaction{}
	}

	// Return the most recent interactions up to the limit
	if limit > 0 && len(interactions) > limit {
		return interactions[len(interactions)-limit:]
	}

	return interactions
}

// GetUserContext retrieves recent conversation context for a user
func (d *Database) GetUserContext(userID int64, messageCount int) string {
	interactions := d.GetUserInteractions(userID, messageCount)
	if len(interactions) == 0 {
		return ""
	}

	var context string
	for _, interaction := range interactions {
		context += fmt.Sprintf("User: %s\nAI: %s\n\n", interaction.UserMessage, interaction.AIResponse)
	}

	return context
}

// GetUserStats retrieves user interaction statistics
func (d *Database) GetUserStats(userID int64) map[string]interface{} {
	interactions := d.GetUserInteractions(userID, 0)
	stats := map[string]interface{}{
		"total_interactions": len(interactions),
		"first_interaction":  nil,
		"last_interaction":   nil,
		"actions_used":       make(map[string]int),
	}

	if len(interactions) > 0 {
		stats["first_interaction"] = interactions[0].Timestamp
		stats["last_interaction"] = interactions[len(interactions)-1].Timestamp

		// Count actions
		for _, interaction := range interactions {
			if interaction.Action != "" {
				stats["actions_used"].(map[string]int)[interaction.Action]++
			}
		}
	}

	return stats
}

// loadInteractions loads interactions from disk
func (d *Database) loadInteractions() error {
	data, err := os.ReadFile(d.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, start with empty database
		}
		return fmt.Errorf("failed to read interactions file: %v", err)
	}

	if len(data) == 0 {
		return nil
	}

	var interactions map[int64][]types.Interaction
	if err := json.Unmarshal(data, &interactions); err != nil {
		return fmt.Errorf("failed to unmarshal interactions: %v", err)
	}

	d.interactions = interactions
	return nil
}

// saveInteractions saves interactions to disk
func (d *Database) saveInteractions() error {
	// Note: This function is called from functions that already hold the write lock
	// so we don't need to acquire any additional locks here
	data, err := json.MarshalIndent(d.interactions, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal interactions: %v", err)
	}

	if err := os.WriteFile(d.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write interactions file: %v", err)
	}

	return nil
}

// Backup creates a backup of the current database
func (d *Database) Backup(backupPath string) error {
	d.mutex.RLock()
	data, err := json.MarshalIndent(d.interactions, "", "  ")
	d.mutex.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal interactions for backup: %v", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}

	return nil
}

// Cleanup removes old interactions (older than specified days)
func (d *Database) Cleanup(daysOld int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	cutoff := time.Now().AddDate(0, 0, -daysOld)
	totalRemoved := 0

	for userID, interactions := range d.interactions {
		var validInteractions []types.Interaction
		for _, interaction := range interactions {
			if interaction.Timestamp.After(cutoff) {
				validInteractions = append(validInteractions, interaction)
			} else {
				totalRemoved++
			}
		}
		d.interactions[userID] = validInteractions
	}

	log.Printf("Cleaned up %d old interactions", totalRemoved)
	return d.saveInteractions()
}
