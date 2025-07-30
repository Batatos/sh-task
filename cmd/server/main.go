package main

import (
	"log"
	"os"

	"skyhawk-security-microservice/internal/database"
	"skyhawk-security-microservice/internal/server"
)

func main() {
	// Connect to database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get port from environment or use default
	port := 8080
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = 8080 // You can parse envPort to int if needed
	}

	// Create and start server
	srv := server.NewServer(db)
	if err := srv.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 