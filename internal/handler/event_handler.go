package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"skyhawk-security-microservice/internal/models"
	"skyhawk-security-microservice/internal/queue"
	"skyhawk-security-microservice/internal/repository"
)

// EventHandler handles security event-related endpoints
type EventHandler struct {
	eventRepo    *repository.EventRepository
	queueManager queue.QueueInterface
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventRepo *repository.EventRepository, queueManager queue.QueueInterface) *EventHandler {
	return &EventHandler{
		eventRepo:    eventRepo,
		queueManager: queueManager,
	}
}

// CreateEvent handles security event creation
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req models.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Create event model
	event := &models.Event{
		EventID:     generateEventID(),
		EventType:   req.EventType,
		Severity:    req.Severity,
		Source:      req.Source,
		Description: req.Description,
		EventData:   req.EventData,
	}

	// Save to database
	if err := h.eventRepo.CreateEvent(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create event",
		})
		return
	}

			// Publish to queue for async processing
		if h.queueManager != nil {
			go func() {
				if err := h.queueManager.PublishEvent(event, "security_events"); err != nil {
					log.Printf("Failed to publish event to queue: %v", err)
				} else {
					log.Printf("Event %s published to queue", event.EventID)
				}
			}()
		}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully and queued for processing",
		"event":   event,
	})
}

// GetEvents handles event retrieval
func (h *EventHandler) GetEvents(c *gin.Context) {
	events, err := h.eventRepo.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve events",
		})
		return
	}

	// Get queue statistics if queue manager is available
	var queueStats map[string]interface{}
	if h.queueManager != nil {
		queueStats = h.queueManager.GetQueueStats("security_events", "security_events_retry", "security_events_dead")
	}

	c.JSON(http.StatusOK, gin.H{
		"events":     events,
		"total":      len(events),
		"queue_stats": queueStats,
	})
}

// GetEvent handles single event retrieval
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")
	
	event, err := h.eventRepo.GetEventByID(eventID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Event not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"event": event,
	})
}

// UpdateEvent handles event updates
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")
	
	var req models.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	event, err := h.eventRepo.UpdateEvent(eventID, &req)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Event not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated successfully",
		"event":   event,
	})
}

// DeleteEvent handles event deletion
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	
	err := h.eventRepo.DeleteEvent(eventID)
	if err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Event not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Event deleted successfully",
		"event_id": eventID,
	})
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return "event-" + time.Now().Format("20060102150405") + "-" + time.Now().Format("000000000")
}

// GetQueueStats handles queue statistics requests
func (h *EventHandler) GetQueueStats(c *gin.Context) {
	if h.queueManager == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Queue manager not available",
		})
		return
	}

	stats := h.queueManager.GetQueueStats("security_events", "security_events_retry", "security_events_dead")
	
	c.JSON(http.StatusOK, gin.H{
		"queue_stats": stats,
		"timestamp":   time.Now(),
	})
} 