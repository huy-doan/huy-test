// backend/internal/middleware/permission.go
package middleware

import (
	"net/http"

	object "github.com/huydq/test/internal/domain/object/permission"
	"github.com/labstack/echo/v4"
)

func (m *MiddlewareManager) RoutePermissions(permissions ...object.PermissionCode) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roleIDValue := c.Get(string(ContextKey_AuthRoleID))
			if roleIDValue == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: missing role information")
			}

			roleID, ok := roleIDValue.(int)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal error: invalid role ID type")
			}

			hasPermission, err := m.permissionMiddlewareService.HasPermission(c.Request().Context(), roleID, permissions...)
			if err != nil {
				m.logger.Error("Error checking permissions", map[string]any{
					"roleID": roleID,
					"error":  err.Error(),
				})
				return echo.NewHTTPError(http.StatusInternalServerError, "Error checking permissions")
			}

			if !hasPermission {
				userPermissions, _ := m.permissionMiddlewareService.GetUserPermissions(c.Request().Context(), roleID)

				requiredPerms := make([]string, len(permissions))
				for i, p := range permissions {
					requiredPerms[i] = string(p)
				}

				m.logger.Warn("Permission denied", map[string]any{
					"roleID":              roleID,
					"requestPath":         c.Request().URL.Path,
					"requiredPermissions": requiredPerms,
					"userPermissions":     userPermissions,
					"userId":              c.Get(string(ContextKey_AuthUserIDKey)),
				})
				return echo.NewHTTPError(http.StatusForbidden, "Forbidden: you don't have permission to access this resource")
			}

			return next(c)
		}
	}
}
