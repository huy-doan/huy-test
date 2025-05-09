// src/api/http/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/huydq/test/src/api/http/response"
	"github.com/huydq/test/src/infrastructure/auth"
)

const (
	UserIDKey   contextKey = "userId"
	EmailKey    contextKey = "email"
	RoleIDKey   contextKey = "roleId"
	RoleCodeKey contextKey = "roleCode"
	TokenKey    contextKey = "token"
)

// NewAuthMiddleware creates middleware for JWT authentication
func NewAuthMiddleware(jwtService *auth.JWTService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "Authorization header is required")
				return
			}

			// Check if header has the correct format
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				response.Unauthorized(w, "Authorization header format must be Bearer {token}")
				return
			}

			// Parse and validate the token
			tokenString := headerParts[1]

			// Check if token is blacklisted
			if jwtService.IsBlacklisted(tokenString) {
				response.Unauthorized(w, "Token has been revoked")
				return
			}

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				response.Unauthorized(w, "Invalid or expired token")
				return
			}

			// Set user ID and role information in the context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)
			ctx = context.WithValue(ctx, RoleIDKey, claims.RoleID)
			ctx = context.WithValue(ctx, RoleCodeKey, claims.RoleCode)
			ctx = context.WithValue(ctx, TokenKey, tokenString) // save token in context for logout

			// Call the next handler with modified context
			next(w, r.WithContext(ctx))
		}
	}
}

// NewRoleMiddleware creates middleware for role-based authorization using role code
func NewRoleMiddleware(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get user role from context (set by AuthMiddleware)
			roleCode, ok := r.Context().Value(RoleCodeKey).(string)
			if !ok {
				response.Unauthorized(w, "Unauthorized: missing role information")
				return
			}

			// Check if user has one of the required roles
			authorized := false
			for _, role := range roles {
				if roleCode == role {
					authorized = true
					break
				}
			}

			if !authorized {
				response.Forbidden(w, "Forbidden: insufficient permissions")
				return
			}

			// Call the next handler
			next(w, r)
		}
	}
}
