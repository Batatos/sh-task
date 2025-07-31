package routes

import (
	"github.com/gin-gonic/gin"
	"skyhawk-security-microservice/internal/handler"
	"skyhawk-security-microservice/internal/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, handlers *handler.Handler) {
	// Apply global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())

	// Health check endpoints
	router.GET("/health", handlers.HealthHandler.HealthCheck)
	router.GET("/", handlers.HealthHandler.GetRoot)
	router.GET("/api/v1/status", handlers.HealthHandler.GetStatus)

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Event routes
		events := apiV1.Group("/events")
		{
			events.POST("/", handlers.EventHandler.CreateEvent)
			events.GET("/", handlers.EventHandler.GetEvents)
			events.GET("/:id", handlers.EventHandler.GetEvent)
			events.PUT("/:id", handlers.EventHandler.UpdateEvent)
			events.DELETE("/:id", handlers.EventHandler.DeleteEvent)
		}

		// Queue routes
		queue := apiV1.Group("/queue")
		{
			queue.GET("/stats", handlers.EventHandler.GetQueueStats)
		}

		// Future route groups can be added here:
		// users := apiV1.Group("/users")
		// incidents := apiV1.Group("/incidents")
		// rules := apiV1.Group("/rules")
	}
} 