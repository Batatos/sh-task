package repository

import (
	"database/sql"
	"fmt"

	"skyhawk-security-microservice/internal/database"
	"skyhawk-security-microservice/internal/models"
)

// EventRepository handles database operations for events
type EventRepository struct {
	db *database.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(db *database.DB) *EventRepository {
	return &EventRepository{db: db}
}

// CreateEvent creates a new event in the database
func (r *EventRepository) CreateEvent(event *models.Event) error {
	query := `
		INSERT INTO security_events (event_id, event_type, severity, source, description, event_data)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		event.EventID,
		event.EventType,
		event.Severity,
		event.Source,
		event.Description,
		event.EventData,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create event: %v", err)
	}

	return nil
}

// GetEventByID retrieves an event by its ID
func (r *EventRepository) GetEventByID(id string) (*models.Event, error) {
	query := `
		SELECT id, event_id, event_type, severity, source, description, event_data, created_at, updated_at
		FROM security_events
		WHERE event_id = $1`

	event := &models.Event{}
	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.EventID,
		&event.EventType,
		&event.Severity,
		&event.Source,
		&event.Description,
		&event.EventData,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %v", err)
	}

	return event, nil
}

// GetAllEvents retrieves all events from the database
func (r *EventRepository) GetAllEvents() ([]*models.Event, error) {
	query := `
		SELECT id, event_id, event_type, severity, source, description, event_data, created_at, updated_at
		FROM security_events
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %v", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		event := &models.Event{}
		err := rows.Scan(
			&event.ID,
			&event.EventID,
			&event.EventType,
			&event.Severity,
			&event.Source,
			&event.Description,
			&event.EventData,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %v", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %v", err)
	}

	return events, nil
}

// UpdateEvent updates an event in the database
func (r *EventRepository) UpdateEvent(eventID string, updates *models.UpdateEventRequest) (*models.Event, error) {
	query := `
		UPDATE security_events
		SET event_type = COALESCE($2, event_type),
			severity = COALESCE($3, severity),
			source = COALESCE($4, source),
			description = COALESCE($5, description),
			event_data = COALESCE($6, event_data),
			updated_at = NOW()
		WHERE event_id = $1
		RETURNING id, event_id, event_type, severity, source, description, event_data, created_at, updated_at`

	event := &models.Event{}
	err := r.db.QueryRow(
		query,
		eventID,
		updates.EventType,
		updates.Severity,
		updates.Source,
		updates.Description,
		updates.EventData,
	).Scan(
		&event.ID,
		&event.EventID,
		&event.EventType,
		&event.Severity,
		&event.Source,
		&event.Description,
		&event.EventData,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to update event: %v", err)
	}

	return event, nil
}

// DeleteEvent deletes an event from the database
func (r *EventRepository) DeleteEvent(eventID string) error {
	query := `DELETE FROM security_events WHERE event_id = $1`

	result, err := r.db.Exec(query, eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found")
	}

	return nil
} 