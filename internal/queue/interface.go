package queue

import (
	"fmt"
	"time"
	"skyhawk-security-microservice/internal/models"
)

// Message represents a message in the queue
type Message struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Retries   int                    `json:"retries"`
}

// QueueInterface defines the interface for queue implementations
type QueueInterface interface {
	PublishMessage(message Message, queueName string) error
	PublishEvent(event *models.Event, queueName string) error
	ConsumeMessage(queueName string, timeout time.Duration) (*Message, error)
	GetQueueLength(queueName string) (int64, error)
	GetQueueStats(queueNames ...string) map[string]interface{}
	Close() error
}

// QueueType represents different queue implementations
type QueueType string

const (
	QueueTypeRabbitMQ QueueType = "rabbitmq"
)

// NewQueue creates a new queue based on the specified type
func NewQueue(queueType QueueType, config map[string]string) (QueueInterface, error) {
	switch queueType {
	case QueueTypeRabbitMQ:
		amqpURL := config["amqp_url"]
		if amqpURL == "" {
			amqpURL = "amqp://admin:password@localhost:5672/"
		}
		return NewRabbitMQQueue(amqpURL)
	
	default:
		return nil, fmt.Errorf("unknown queue type: %s", queueType)
	}
} 