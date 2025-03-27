// src/api/graphql/middleware/logger.go

package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
)

const (
	LoggerContextKey = "graphql_logger"
)

// GraphQLLoggerMiddleware creates middleware for adding logger to GraphQL context
func GraphQLLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the global logger instance
		appLogger := logger.GetLogger()
		
		// Get trace ID from request header or generate a new one
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = logger.GenerateTraceID()
		}
		
		// Set trace ID in response header for debugging
		c.Header("X-Trace-ID", traceID)
		
		// Create request-specific logger with trace ID
		reqLogger := appLogger.WithTraceID(traceID)

		// Store logger in Gin context
		c.Set("logger", reqLogger)
		
		// Log basic request info
		reqLogger.Info("GraphQL request received", map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})
		
		c.Next()
	}
}

// WithLogger adds a logger to GraphQL context
func WithLogger(ctx context.Context, l logger.Logger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, l)
}

// GetLogger retrieves the logger from GraphQL context
func GetLogger(ctx context.Context) logger.Logger {
	if l, ok := ctx.Value(LoggerContextKey).(logger.Logger); ok {
		return l
	}
	// Return the global logger instance if none is found in context
	return logger.GetLogger()
}
