package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"skyhawk-security-microservice/internal/database"
	"skyhawk-security-microservice/internal/handler"
)

func TestSetupRoutes(t *testing.T) {
	// Create a mock database connection
	mockDB := &database.DB{}

	// Create handlers
	handlers := handler.NewHandler(mockDB)

	// Create router
	router := gin.New()
	SetupRoutes(router, handlers)

	// Test health endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAPIRoutes(t *testing.T) {
	// Create a mock database connection
	mockDB := &database.DB{}

	// Create handlers
	handlers := handler.NewHandler(mockDB)

	// Create router
	router := gin.New()
	SetupRoutes(router, handlers)

	// Test API status endpoint
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/status", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
} 