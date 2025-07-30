package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection() (*DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "skyhawk_security")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection with retries
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err != nil {
			log.Printf("Attempt %d: Failed to ping database: %v", i+1, err)
			if i < maxRetries-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to ping database after %d attempts: %v", maxRetries, err)
		}
		break
	}

	log.Println("✅ Successfully connected to PostgreSQL database")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
} 