package middleware

import (
	"net/http"
	"time"

	"github.com/huydq/test/src/infrastructure/logger"
)

// PerformanceThreshold is the duration in milliseconds above which a request is considered slow
const PerformanceThreshold = 500 // milliseconds

// PerformanceMonitor logs performance metrics for API requests
func PerformanceMonitor(appLogger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start timing
			startTime := time.Now()

			// Get trace ID if it exists, or use a new one
			traceID := r.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = logger.GenerateTraceID()
			}

			// Create a logger with the trace ID
			perfLogger := appLogger.WithTraceID(traceID)

			// Process request
			next.ServeHTTP(w, r)

			// Calculate duration
			duration := time.Since(startTime)

			// Log performance metrics for slow requests only or when in debug mode
			// This avoids duplicate logs for normal requests
			if duration.Milliseconds() > PerformanceThreshold {
				fields := map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"duration_ms": duration.Milliseconds(),
					"is_slow":     true,
				}
				perfLogger.Warn("Slow API request", fields)
			}
		})
	}
}
