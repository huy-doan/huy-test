package middleware

import (
	"github.com/huydq/test/internal/pkg/database"
	"github.com/huydq/test/internal/pkg/logger"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ContextKey is used to prevent key collisions in context
type ContextKey string

// DBMiddleware creates a middleware that adds the database connection to the context
func (m *MiddlewareManager) DBContextMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, err := database.SetDB(c.Request().Context(), db)
			if err != nil {
				logger.GetLogger().Error("Failed to set DB in context", map[string]any{
					"error": err,
				})
				return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Database error")
			}

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
