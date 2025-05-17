package middleware

import (
	"time"

	"github.com/huydq/test/internal/pkg/logger"
	"github.com/labstack/echo/v4"
)

const (
	PerformanceThreshold = 500 // milliseconds
)

func (m *MiddlewareManager) PerformanceMonitor(appLogger logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Start timing
			startTime := time.Now()

			// Get trace ID if it exists, or use a new one
			traceID := c.Request().Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = logger.GenerateTraceID()
				c.Request().Header.Set("X-Trace-ID", traceID)
				c.Response().Header().Set("X-Trace-ID", traceID)
			}

			// Create a logger with the trace ID
			perfLogger := appLogger.WithTraceID(traceID)

			// Set trace ID in context for other middleware/handlers
			c.Set("traceID", traceID)

			// Process request
			err := next(c)

			// Calculate duration
			duration := time.Since(startTime)

			// Log performance metrics for slow requests
			if duration.Milliseconds() > PerformanceThreshold {
				fields := map[string]any{
					"method":      c.Request().Method,
					"path":        c.Request().URL.Path,
					"status_code": c.Response().Status,
					"duration_ms": duration.Milliseconds(),
					"is_slow":     true,
				}
				perfLogger.Warn("Slow API request", fields)
			}

			return err
		}
	}
}
