package middleware

import (
	"net/http"
	"strings"

	tokenDomainSvc "github.com/huydq/test/internal/domain/service/auth"
	"github.com/huydq/test/internal/infrastructure/adapter/auth"
	messages "github.com/huydq/test/internal/pkg/utils/messages"
	"github.com/labstack/echo/v4"
)

const (
	ContextKey_AuthUserIDKey ContextKey = "userId"
	ContextKey_AuthEmail     ContextKey = "email"
	ContextKey_AuthRoleID    ContextKey = "roleId"
	ContextKey_AuthToken     ContextKey = "token"
)

// JWT creates middleware for JWT authentication
func (m *MiddlewareManager) JWTMiddleware(jwtService *auth.JWTService, tokenDomainSvc tokenDomainSvc.AccessTokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, messages.MsgUnauthorized)
			}

			// Check if header has the correct format
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, messages.MsgInvalidTokenFormat)
			}

			// Parse and validate the token
			tokenString := headerParts[1]
			isBlacklisted, err := tokenDomainSvc.IsBlacklisted(c.Request().Context(), tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, messages.MsgInvalidToken)
			}

			if isBlacklisted {
				return echo.NewHTTPError(http.StatusUnauthorized, messages.MsgTokenBlacklisted)
			}

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, messages.MsgInvalidToken)
			}

			// Set user ID and role information in the context
			c.Set(string(ContextKey_AuthUserIDKey), claims.UserID)
			c.Set(string(ContextKey_AuthEmail), claims.Email)
			c.Set(string(ContextKey_AuthRoleID), claims.RoleID)
			c.Set(string(ContextKey_AuthToken), tokenString)

			// Call the next handler
			return next(c)
		}
	}
}
