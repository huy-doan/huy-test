package middleware

import (
	"net/http"
	"time"

	"github.com/huydq/test/src/infrastructure/logger"
)

// RequestLoggerMiddleware logs information about incoming HTTP requests
func RequestLoggerMiddleware(appLogger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate a unique trace ID for this request if not already set
			traceID := r.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = logger.GenerateTraceID()
				r.Header.Set("X-Trace-ID", traceID)
			}

			requestLogger := appLogger.WithTraceID(traceID)

			// Log the request
			startTime := time.Now()

			// Create a custom response writer to capture the status code
			rw := &responseStatusWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // Default to 200 OK
			}

			// Add trace ID to response header for debugging and tracking
			rw.Header().Set("X-Trace-ID", traceID)

			// Process the request
			defer func() {
				if err := recover(); err != nil {
					// Log panic with stack trace
					responseFields := map[string]interface{}{
						"status_code": http.StatusInternalServerError,
						"duration_ms": time.Since(startTime).Milliseconds(),
						"error":       err,
					}
					requestLogger.Error("Request panic", responseFields)

					// Return 500 Internal Server Error
					http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			// Call the next handler
			next.ServeHTTP(rw, r)
		})
	}
}

// responseStatusWriter is a wrapper for http.ResponseWriter that captures the status code
type responseStatusWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and passes it to the wrapped ResponseWriter
func (rw *responseStatusWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// getClientIP extracts the client IP from various headers or from the request
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// The client IP is the first IP in the list
		for i := 0; i < len(forwardedFor); i++ {
			if forwardedFor[i] == ',' {
				return forwardedFor[:i]
			}
		}
		return forwardedFor
	}

	// Check for X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
