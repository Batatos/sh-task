package handler

import (
	"log"
	"os"
	"skyhawk-security-microservice/internal/database"
	"skyhawk-security-microservice/internal/queue"
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

		// Create RabbitMQ queue manager
	var queueManager queue.QueueInterface
	
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://admin:password@rabbitmq:5672/"
	}

	var err error
	queueManager, err = queue.NewRabbitMQQueue(amqpURL)
	if err != nil {
		log.Printf("Warning: Failed to create RabbitMQ queue manager: %v", err)
		log.Printf("Queue functionality will be disabled")
		queueManager = nil
	} else {
		log.Printf("RabbitMQ queue manager initialized successfully")
	}

	return &Handler{
		HealthHandler: NewHealthHandler(),
		EventHandler:  NewEventHandler(eventRepo, queueManager),
	}
} 