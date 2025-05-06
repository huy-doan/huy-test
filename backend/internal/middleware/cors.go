package middleware

import (
	"net/http"
	"os"
	"slices"
	"strings"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		allowMethods := "GET,POST,PUT,PATCH,DELETE,OPTIONS"
		if os.Getenv("CORS_ALLOW_METHODS") != "" {
			allowMethods = os.Getenv("CORS_ALLOW_METHODS")
		}

		// Get allowed headers from environment, default to standard headers
		allowHeaders := "Origin,Content-Type,Accept,Authorization"
		if os.Getenv("CORS_ALLOW_HEADERS") != "" {
			allowHeaders = os.Getenv("CORS_ALLOW_HEADERS")
		}

		// Get the origin from the request
		origin := r.Header.Get("Origin")

		// Check if the origin is allowed
		allowOrigin := "*"
		if origin != "" && len(allowedOrigins) > 0 {
			if allowedOrigins[0] != "*" {
				allowOrigin = ""
				// Check if the request origin is in the allowed origins list
				if slices.Contains(allowedOrigins, origin) {
					allowOrigin = origin
				}
			}
		}

		// Set CORS headers
		if allowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", allowMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
