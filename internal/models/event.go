package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Event represents a security event
type Event struct {
	ID          string    `json:"id" db:"id"`
	EventID     string    `json:"event_id" db:"event_id"`
	EventType   string    `json:"event_type" db:"event_type"`
	Severity    string    `json:"severity" db:"severity"`
	Source      string    `json:"source" db:"source"`
	Description string    `json:"description" db:"description"`
	EventData   EventData `json:"event_data" db:"event_data"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// EventData represents the JSON data for an event
type EventData map[string]interface{}

// Value implements the driver.Valuer interface for JSONB
func (e EventData) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	return json.Marshal(e)
}

// Scan implements the sql.Scanner interface for JSONB
func (e *EventData) Scan(value interface{}) error {
	if value == nil {
		*e = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, e)
}

// CreateEventRequest represents the request to create an event
type CreateEventRequest struct {
	EventType   string    `json:"event_type" binding:"required"`
	Severity    string    `json:"severity" binding:"required"`
	Source      string    `json:"source" binding:"required"`
	Description string    `json:"description"`
	EventData   EventData `json:"event_data"`
}

// UpdateEventRequest represents the request to update an event
type UpdateEventRequest struct {
	EventType   string    `json:"event_type"`
	Severity    string    `json:"severity"`
	Source      string    `json:"source"`
	Description string    `json:"description"`
	EventData   EventData `json:"event_data"`
} 