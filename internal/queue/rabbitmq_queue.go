package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"skyhawk-security-microservice/internal/models"
)

// RabbitMQQueue implements queue using RabbitMQ
type RabbitMQQueue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewRabbitMQQueue creates a new RabbitMQ queue manager
func NewRabbitMQQueue(amqpURL string) (*RabbitMQQueue, error) {
	// Parse AMQP URL
	if amqpURL == "" {
		amqpURL = "amqp://admin:password@localhost:5672/"
	}

	var conn *amqp.Connection
	var err error
	
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			break
		}
		
		log.Printf("Attempt %d: Failed to connect to RabbitMQ: %v", i+1, err)
		if i < maxRetries-1 {
			time.Sleep(2 * time.Second)
			continue
		}
		return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
	}

	// Create channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	queue := &RabbitMQQueue{
		conn:    conn,
		channel: channel,
		ctx:     ctx,
		cancel:  cancel,
	}

	log.Printf("Connected to RabbitMQ successfully")
	return queue, nil
}

// PublishMessage publishes a message to a queue
func (rq *RabbitMQQueue) PublishMessage(message Message, queueName string) error {
	// Declare queue
	_, err := rq.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Serialize message
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish message
	err = rq.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         messageBytes,
			DeliveryMode: amqp.Persistent, // Make message persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message %s to RabbitMQ queue %s", message.ID, queueName)
	return nil
}

// PublishEvent publishes an event to the queue
func (rq *RabbitMQQueue) PublishEvent(event *models.Event, queueName string) error {
	message := Message{
		ID:        event.EventID,
		Type:      "security_event",
		Data:      map[string]interface{}{"event": event},
		Timestamp: time.Now(),
		Retries:   0,
	}

	return rq.PublishMessage(message, queueName)
}

// ConsumeMessage consumes a message from a queue
func (rq *RabbitMQQueue) ConsumeMessage(queueName string, timeout time.Duration) (*Message, error) {
	// Declare queue
	_, err := rq.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Set QoS for fair dispatch
	err = rq.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	// Consume messages
	msgs, err := rq.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consuming: %w", err)
	}

	// Wait for message with timeout
	select {
	case msg := <-msgs:
		// Parse message
		var message Message
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			msg.Nack(false, true) // Reject and requeue
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}

		// Acknowledge message
		msg.Ack(false)

		log.Printf("Consumed message %s from RabbitMQ queue %s", message.ID, queueName)
		return &message, nil

	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for message")

	case <-rq.ctx.Done():
		return nil, fmt.Errorf("queue context cancelled")
	}
}

// StartConsumer starts a consumer that continuously processes messages
func (rq *RabbitMQQueue) StartConsumer(queueName string, workerID int) {
	log.Printf("Starting RabbitMQ consumer worker %d for queue %s", workerID, queueName)

	// Declare queue
	_, err := rq.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return
	}

	// Set QoS for fair dispatch
	err = rq.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Printf("Failed to set QoS: %v", err)
		return
	}

	// Consume messages
	msgs, err := rq.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Printf("Failed to start consuming: %v", err)
		return
	}

	// Process messages
	for {
		select {
		case msg := <-msgs:
			// Parse message
			var message Message
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, true) // Reject and requeue
				continue
			}

			// Process the message
			if err := rq.ProcessEvent(&message); err != nil {
				log.Printf("Error processing message %s: %v", message.ID, err)
				
				// Increment retry count
				message.Retries++
				
				// If max retries not reached, requeue
				if message.Retries < 3 {
					log.Printf("Requeuing message %s (retry %d)", message.ID, message.Retries)
					if err := rq.PublishMessage(message, queueName+"_retry"); err != nil {
						log.Printf("Failed to requeue message: %v", err)
					}
					msg.Ack(false) // Acknowledge original message
				} else {
					log.Printf("Message %s exceeded max retries, moving to dead letter queue", message.ID)
					if err := rq.PublishMessage(message, queueName+"_dead"); err != nil {
						log.Printf("Failed to move message to dead letter queue: %v", err)
					}
					msg.Ack(false) // Acknowledge original message
				}
			} else {
				// Successfully processed
				msg.Ack(false)
			}

		case <-rq.ctx.Done():
			log.Printf("Consumer worker %d stopping", workerID)
			return
		}
	}
}

// ProcessEvent processes a security event message (same as Redis implementation)
func (rq *RabbitMQQueue) ProcessEvent(message *Message) error {
	log.Printf("Processing event: %s", message.ID)

	// Extract event data
	eventData, ok := message.Data["event"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data in message")
	}

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Simulate different processing based on event type
	eventType, _ := eventData["event_type"].(string)
	switch eventType {
	case "login":
		log.Printf("Processing login event: %s", message.ID)
		// Simulate login processing
		time.Sleep(50 * time.Millisecond)
	case "data_access":
		log.Printf("Processing data access event: %s", message.ID)
		// Simulate data access processing
		time.Sleep(75 * time.Millisecond)
	case "file_access":
		log.Printf("Processing file access event: %s", message.ID)
		// Simulate file access processing
		time.Sleep(60 * time.Millisecond)
	default:
		log.Printf("Processing generic event: %s", message.ID)
	}

	log.Printf("Successfully processed event: %s", message.ID)
	return nil
}

// GetQueueLength returns the number of messages in a queue
func (rq *RabbitMQQueue) GetQueueLength(queueName string) (int64, error) {
	// Declare queue to get info
	queue, err := rq.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return 0, fmt.Errorf("failed to declare queue: %w", err)
	}

	return int64(queue.Messages), nil
}

// GetQueueStats returns statistics about queues
func (rq *RabbitMQQueue) GetQueueStats(queueNames ...string) map[string]interface{} {
	stats := make(map[string]interface{})

	for _, queueName := range queueNames {
		length, err := rq.GetQueueLength(queueName)
		if err != nil {
			stats[queueName] = map[string]interface{}{
				"error": err.Error(),
			}
			continue
		}

		stats[queueName] = map[string]interface{}{
			"length": length,
			"type":   "rabbitmq",
		}
	}

	return stats
}

// Close closes the RabbitMQ connection
func (rq *RabbitMQQueue) Close() error {
	rq.cancel()
	
	if rq.channel != nil {
		rq.channel.Close()
	}
	
	if rq.conn != nil {
		return rq.conn.Close()
	}
	
	return nil
} 