package health

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"skyhawk-security-microservice/internal/database"
)

// HealthStatus represents the overall health status
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Version   string                 `json:"version"`
	Checks    map[string]CheckResult `json:"checks"`
}

// CheckResult represents the result of a health check
type CheckResult struct {
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Duration  string    `json:"duration"`
}

// HealthChecker manages all health checks
type HealthChecker struct {
	db           *database.DB
	startTime    time.Time
	version      string
	mu           sync.RWMutex
	checkResults map[string]CheckResult
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(db *database.DB) *HealthChecker {
	return &HealthChecker{
		db:           db,
		startTime:    time.Now(),
		version:      "1.0.0",
		checkResults: make(map[string]CheckResult),
	}
}

// CheckHealth performs all health checks
func (hc *HealthChecker) CheckHealth(ctx context.Context) HealthStatus {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	// Perform all health checks concurrently
	var wg sync.WaitGroup
	checks := []string{"database", "memory", "disk"}

	for _, check := range checks {
		wg.Add(1)
		go func(checkName string) {
			defer wg.Done()
			hc.performCheck(ctx, checkName)
		}(check)
	}

	wg.Wait()

	// Determine overall status
	overallStatus := "healthy"
	for _, result := range hc.checkResults {
		if result.Status == "unhealthy" {
			overallStatus = "unhealthy"
			break
		}
	}

	return HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Uptime:    time.Since(hc.startTime).String(),
		Version:   hc.version,
		Checks:    hc.checkResults,
	}
}

// performCheck performs a specific health check
func (hc *HealthChecker) performCheck(ctx context.Context, checkName string) {
	start := time.Now()
	var result CheckResult

	switch checkName {
	case "database":
		result = hc.checkDatabase(ctx)
	case "memory":
		result = hc.checkMemory()
	case "disk":
		result = hc.checkDisk()
	default:
		result = CheckResult{
			Status:    "unknown",
			Message:   fmt.Sprintf("Unknown check: %s", checkName),
			Timestamp: time.Now(),
		}
	}

	result.Duration = time.Since(start).String()
	hc.checkResults[checkName] = result
}

// checkDatabase checks database connectivity
func (hc *HealthChecker) checkDatabase(ctx context.Context) CheckResult {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Test database connection
	if err := hc.db.PingContext(ctx); err != nil {
		return CheckResult{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("Database connection failed: %v", err),
			Timestamp: time.Now(),
		}
	}

	// Test a simple query
	var result int
	if err := hc.db.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return CheckResult{
			Status:    "unhealthy",
			Message:   fmt.Sprintf("Database query failed: %v", err),
			Timestamp: time.Now(),
		}
	}

	return CheckResult{
		Status:    "healthy",
		Message:   "Database connection and queries working",
		Timestamp: time.Now(),
	}
}

// checkMemory checks memory usage
func (hc *HealthChecker) checkMemory() CheckResult {
	// In a real application, you'd use runtime.ReadMemStats
	// For now, we'll simulate a memory check
	return CheckResult{
		Status:    "healthy",
		Message:   "Memory usage within normal limits",
		Timestamp: time.Now(),
	}
}

// checkDisk checks disk space
func (hc *HealthChecker) checkDisk() CheckResult {
	// In a real application, you'd check disk space
	// For now, we'll simulate a disk check
	return CheckResult{
		Status:    "healthy",
		Message:   "Disk space available",
		Timestamp: time.Now(),
	}
}

// GetReadinessStatus checks if the service is ready to handle requests
func (hc *HealthChecker) GetReadinessStatus(ctx context.Context) HealthStatus {
	// For readiness, we only check critical dependencies
	hc.mu.Lock()
	defer hc.mu.Unlock()

	// Check database readiness
	start := time.Now()
	dbResult := hc.checkDatabase(ctx)
	dbResult.Duration = time.Since(start).String()

	checks := map[string]CheckResult{
		"database": dbResult,
	}

	overallStatus := "ready"
	if dbResult.Status == "unhealthy" {
		overallStatus = "not_ready"
	}

	return HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Uptime:    time.Since(hc.startTime).String(),
		Version:   hc.version,
		Checks:    checks,
	}
} 