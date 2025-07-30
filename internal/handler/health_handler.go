package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "skyhawk-security-microservice",
		"version":   "1.0.0",
	})
}

func (h *HealthHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "operational",
		"uptime":    "running",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (h *HealthHandler) GetRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "Skyhawk Security Microservice",
		"version": "1.0.0",
		"status":  "running",
	})
} 