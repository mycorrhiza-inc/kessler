package health

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"kessler/internal/cache"
	"kessler/internal/database"
	"kessler/pkg/logger"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// HealthStatus represents the overall health status
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusDegraded  HealthStatus = "degraded"
	StatusUnhealthy HealthStatus = "unhealthy"
)

// ComponentHealth represents the health of an individual component
type ComponentHealth struct {
	Status      HealthStatus `json:"status"`
	Message     string       `json:"message,omitempty"`
	Duration    string       `json:"duration"`
	LastChecked time.Time    `json:"last_checked"`
}

// HealthResponse represents the complete health check response
type HealthResponse struct {
	Status     HealthStatus               `json:"status"`
	Timestamp  time.Time                  `json:"timestamp"`
	Duration   string                     `json:"duration"`
	Version    string                     `json:"version,omitempty"`
	Components map[string]ComponentHealth `json:"components"`
}

// SimpleHealthResponse for basic liveness checks
type SimpleHealthResponse struct {
	Status    HealthStatus `json:"status"`
	Timestamp time.Time    `json:"timestamp"`
	Message   string       `json:"message"`
}

func DefineHealthRoutes(r *mux.Router) {
	// Kubernetes-style endpoints
	r.HandleFunc("/", HealthHandler).Methods(http.MethodGet)         // Detailed health
	r.HandleFunc("/live", LivenessHandler).Methods(http.MethodGet)   // Liveness probe
	r.HandleFunc("/ready", ReadinessHandler).Methods(http.MethodGet) // Readiness probe

	// Legacy endpoints (keeping your existing ones but making them GET)
	r.HandleFunc("/complete-check", HealthHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/minimal-check", LivenessHandler).Methods(http.MethodGet, http.MethodPost)
}

// LivenessHandler checks if the application is alive (basic check)
func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	start := time.Now()

	response := SimpleHealthResponse{
		Status:    StatusHealthy,
		Timestamp: time.Now(),
		Message:   "Service is alive",
	}

	log.DebugContext(ctx, "Liveness check completed",
		zap.Duration("duration", time.Since(start)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ReadinessHandler checks if the application is ready to serve traffic
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	start := time.Now()

	// Check critical dependencies
	components := make(map[string]ComponentHealth)
	overallStatus := StatusHealthy

	// Check database
	dbHealth := checkDatabase(ctx)
	components["database"] = dbHealth
	if dbHealth.Status != StatusHealthy {
		overallStatus = StatusUnhealthy
	}

	response := HealthResponse{
		Status:     overallStatus,
		Timestamp:  time.Now(),
		Duration:   time.Since(start).String(),
		Components: components,
	}

	statusCode := http.StatusOK
	if overallStatus == StatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	log.InfoContext(ctx, "Readiness check completed",
		zap.String("status", string(overallStatus)),
		zap.Duration("duration", time.Since(start)))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// HealthHandler performs comprehensive health checks
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	start := time.Now()

	components := make(map[string]ComponentHealth)
	overallStatus := StatusHealthy

	// Check database
	dbHealth := checkDatabase(ctx)
	components["database"] = dbHealth
	if dbHealth.Status == StatusUnhealthy {
		overallStatus = StatusUnhealthy
	} else if dbHealth.Status == StatusDegraded && overallStatus == StatusHealthy {
		overallStatus = StatusDegraded
	}

	// Check cache
	cacheHealth := checkCache(ctx)
	components["cache"] = cacheHealth
	if cacheHealth.Status == StatusUnhealthy {
		if overallStatus == StatusHealthy {
			overallStatus = StatusDegraded // Cache failure is degraded, not unhealthy
		}
	}

	// Add version info if available
	version := getVersion()

	response := HealthResponse{
		Status:     overallStatus,
		Timestamp:  time.Now(),
		Duration:   time.Since(start).String(),
		Version:    version,
		Components: components,
	}

	statusCode := http.StatusOK
	if overallStatus == StatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	} else if overallStatus == StatusDegraded {
		statusCode = http.StatusOK // Still serving traffic, but degraded
	}

	log.InfoContext(ctx, "Health check completed",
		zap.String("status", string(overallStatus)),
		zap.Duration("duration", time.Since(start)),
		zap.Int("status_code", statusCode))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// checkDatabase verifies database connectivity
func checkDatabase(ctx context.Context) ComponentHealth {
	start := time.Now()

	// Create a timeout context for the database check
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	q := database.GetTx()
	_, err := q.HealthCheck(checkCtx)

	duration := time.Since(start)

	if err != nil {
		return ComponentHealth{
			Status:      StatusUnhealthy,
			Message:     "Database connection failed: " + err.Error(),
			Duration:    duration.String(),
			LastChecked: time.Now(),
		}
	}

	return ComponentHealth{
		Status:      StatusHealthy,
		Message:     "Database connection successful",
		Duration:    duration.String(),
		LastChecked: time.Now(),
	}
}

// checkCache verifies cache connectivity
func checkCache(ctx context.Context) ComponentHealth {
	start := time.Now()

	// Create a timeout context for the cache check
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Check if cache client exists
	if cache.MemcachedClient == nil {
		return ComponentHealth{
			Status:      StatusDegraded,
			Message:     "Cache client not initialized",
			Duration:    time.Since(start).String(),
			LastChecked: time.Now(),
		}
	}

	// Ping the cache with timeout
	done := make(chan error, 1)
	go func() {
		done <- cache.MemcachedClient.Ping()
	}()

	select {
	case err := <-done:
		duration := time.Since(start)
		if err != nil {
			return ComponentHealth{
				Status:      StatusDegraded,
				Message:     "Cache ping failed: " + err.Error(),
				Duration:    duration.String(),
				LastChecked: time.Now(),
			}
		}

		return ComponentHealth{
			Status:      StatusHealthy,
			Message:     "Cache connection successful",
			Duration:    duration.String(),
			LastChecked: time.Now(),
		}

	case <-checkCtx.Done():
		return ComponentHealth{
			Status:      StatusDegraded,
			Message:     "Cache check timeout",
			Duration:    time.Since(start).String(),
			LastChecked: time.Now(),
		}
	}
}

// getVersion returns the application version
func getVersion() string {
	// You can get this from environment variable or build info
	if version := os.Getenv("VERSION_HASH"); version != "" {
		return version
	}
	return "unknown"
}

// Legacy functions for backward compatibility
func CompleteHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	HealthHandler(w, r)
}

func MinimalHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	LivenessHandler(w, r)
}

func MinimalHealthCheck(r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	q := database.GetTx()
	_, err := q.HealthCheck(ctx)
	return err
}
