package middleware

import (
	"net/http"
	"os"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get allowed origins from environment, default to "*"
		allowOrigin := "*"
		if os.Getenv("CORS_ALLOW_ORIGINS") != "" {
			allowOrigin = os.Getenv("CORS_ALLOW_ORIGINS")
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

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", allowMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
