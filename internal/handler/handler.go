package handler

import (
	"github.com/gin-gonic/gin"
	"skyhawk-security-microservice/internal/database"
	"skyhawk-security-microservice/internal/repository"
)

// Handler coordinates all HTTP handlers
type Handler struct {
	HealthHandler *HealthHandler
	EventHandler  *EventHandler
	// Add more handlers as you add them
	// UserHandler    *UserHandler
	// AuthHandler    *AuthHandler
}

// NewHandler creates a new handler coordinator
func NewHandler(db *database.DB) *Handler {
	eventRepo := repository.NewEventRepository(db)

	return &Handler{
		HealthHandler: NewHealthHandler(),
		EventHandler:  NewEventHandler(eventRepo),
	}
} 