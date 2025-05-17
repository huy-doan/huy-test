package middleware

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (m *MiddlewareManager) CORSMiddleware() echo.MiddlewareFunc {
	// Get allowed origins from environment, default to "*"
	allowedOrigins := []string{"*"}
	if os.Getenv("CORS_ALLOW_ORIGINS") != "" {
		allowedOrigins = strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ",")
		// Trim spaces from each origin
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}

	// Get allowed methods from environment, default to standard methods
	allowMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	if os.Getenv("CORS_ALLOW_METHODS") != "" {
		allowMethods = strings.Split(os.Getenv("CORS_ALLOW_METHODS"), ",")
		// Trim spaces from each method
		for i := range allowMethods {
			allowMethods[i] = strings.TrimSpace(allowMethods[i])
		}
	}

	// Get allowed headers from environment, default to standard headers
	allowHeaders := []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	if os.Getenv("CORS_ALLOW_HEADERS") != "" {
		allowHeaders = strings.Split(os.Getenv("CORS_ALLOW_HEADERS"), ",")
		// Trim spaces from each header
		for i := range allowHeaders {
			allowHeaders[i] = strings.TrimSpace(allowHeaders[i])
		}
	}

	// Configure CORS using Echo's middleware
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	})
}
