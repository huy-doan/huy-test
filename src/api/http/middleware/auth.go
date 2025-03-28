package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vnlab/makeshop-payment/src/api/http/response"
	"github.com/vnlab/makeshop-payment/src/infrastructure/auth"
)

// AuthMiddleware creates middleware for JWT authentication
func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if header has the correct format
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.Unauthorized(c, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		// Parse and validate the token
		tokenString := headerParts[1]

		// Check if token is blacklisted
		if jwtService.IsBlacklisted(tokenString) {
			response.Unauthorized(c, "Token has been revoked")
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user ID and role information in the context
		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("roleId", claims.RoleID)
		c.Set("roleCode", claims.RoleCode)
		c.Set("token", tokenString) // save token in context for logout

		c.Next()
	}
}

// RoleMiddleware creates middleware for role-based authorization using role code
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		roleCode, exists := c.Get("roleCode")
		if !exists {
			response.Unauthorized(c, "Unauthorized: missing role information")
			c.Abort()
			return
		}

		// Check if user has one of the required roles
		userRoleCode, ok := roleCode.(string)
		if !ok {
			response.Unauthorized(c, "Unauthorized: invalid role format")
			c.Abort()
			return
		}

		authorized := false
		for _, role := range roles {
			if userRoleCode == role {
				authorized = true
				break
			}
		}

		if !authorized {
			response.Forbidden(c, "Forbidden: insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
