package middleware

import (
	"net/http"
	"strings"

	"slices"

	"github.com/labstack/echo/v4"

	"github.com/huydq/test/src/infrastructure/auth"
)

const (
	ContextKey_AuthUserIDKey ContextKey = "userId"
	ContextKey_AuthEmail     ContextKey = "email"
	ContextKey_AuthRoleID    ContextKey = "roleId"
	ContextKey_AuthRoleCode  ContextKey = "roleCode"
	ContextKey_AuthToken     ContextKey = "token"
)

// JWT creates middleware for JWT authentication
func (m *MiddlewareManager) JWTMiddleware(jwtService *auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
			}

			// Check if header has the correct format
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			}

			// Parse and validate the token
			tokenString := headerParts[1]

			// Check if token is blacklisted
			if jwtService.IsBlacklisted(tokenString) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token has been revoked")
			}

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			// Set user ID and role information in the context
			c.Set(string(ContextKey_AuthUserIDKey), claims.UserID)
			c.Set(string(ContextKey_AuthEmail), claims.Email)
			c.Set(string(ContextKey_AuthRoleID), claims.RoleID)
			c.Set(string(ContextKey_AuthRoleCode), claims.RoleCode)
			c.Set(string(ContextKey_AuthToken), tokenString) // save token in context for logout

			// Call the next handler
			return next(c)
		}
	}
}

// RoleAuthorization creates middleware for role-based authorization using role code
func (m *MiddlewareManager) RoleAuthorization(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role from context (set by JWT middleware)
			roleCode, ok := c.Get(string(ContextKey_AuthRoleCode)).(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: missing role information")
			}

			// Check if user has one of the required roles
			authorized := slices.Contains(roles, roleCode)

			if !authorized {
				return echo.NewHTTPError(http.StatusForbidden, "Forbidden: insufficient permissions")
			}

			// Call the next handler
			return next(c)
		}
	}
}
