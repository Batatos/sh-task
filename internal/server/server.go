package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"skyhawk-security-microservice/internal/database"
	"skyhawk-security-microservice/internal/handler"
	"skyhawk-security-microservice/internal/routes"
)

type Server struct {
	router *gin.Engine
	server *http.Server
	db     *database.DB
}

func NewServer(db *database.DB) *Server {
	// Set Gin mode
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Setup routes and middleware
	handlers := handler.NewHandler(db)
	routes.SetupRoutes(router, handlers)

	return &Server{
		router: router,
		db:     db,
	}
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      s.router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on :%d", port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
	return nil
}

// GetRouter returns the router for testing purposes
func (s *Server) GetRouter() *gin.Engine {
	return s.router
} 