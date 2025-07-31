package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"skyhawk-security-microservice/internal/queue"
)

func main() {
	// Parse command line flags
	amqpURL := flag.String("amqp", "amqp://admin:password@localhost:5672/", "AMQP URL")
	queueName := flag.String("queue", "security_events", "Queue name")
	workers := flag.Int("workers", 3, "Number of worker goroutines")
	flag.Parse()

	log.Printf("Starting RabbitMQ worker service...")
	log.Printf("AMQP URL: %s", *amqpURL)
	log.Printf("Queue: %s", *queueName)
	log.Printf("Workers: %d", *workers)

	// Create queue manager
	queueManager, err := queue.NewRabbitMQQueue(*amqpURL)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ queue manager: %v", err)
	}
	defer queueManager.Close()

	// Create wait group for workers
	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= *workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			queueManager.StartConsumer(*queueName, workerID)
		}(i)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Queue worker service started. Press Ctrl+C to stop.")

	// Wait for signal
	<-sigChan
	log.Printf("Shutting down queue worker service...")

	// Wait for all workers to finish
	wg.Wait()
	log.Printf("Queue worker service stopped.")
} 