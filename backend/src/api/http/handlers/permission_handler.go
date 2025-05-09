package handlers

import (
	"net/http"

	"github.com/huydq/test/src/api/http/errors"
	"github.com/huydq/test/src/api/http/middleware"
	"github.com/huydq/test/src/api/http/response"
	"github.com/huydq/test/src/api/http/serializers"
	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/infrastructure/logger"
	"github.com/huydq/test/src/lib/i18n"
	"github.com/huydq/test/src/usecase"
)

type PermissionHandler struct {
	permissionUsecase *usecase.PermissionUsecase
}

func NewPermissionHandler(permissionUsecase *usecase.PermissionUsecase) *PermissionHandler {
	return &PermissionHandler{
		permissionUsecase: permissionUsecase,
	}
}

// ListPermission handles the request to list permissions
// @Summary List all permissions
// @Description Get a list of all permissions in the system with their associated screens
// @Tags Admin Permission Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=map[string][]serializers.PermissionResponse} "Success"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /api/v1/admin/permissions [get]
func (h *PermissionHandler) ListPermission(w http.ResponseWriter, r *http.Request) {
	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	permissions, err := h.permissionUsecase.ListPermission(r.Context())
	if err != nil {
		logger.GetLogger().Error("Failed to list permissions", map[string]any{"error": err.Error()})
		response.Error(w, errors.InternalError(i18n.T(r.Context(), "common.error")))
		return
	}

	responseData := map[string]any{
		"permissions": serializers.SerializePermissionCollection(permissions),
	}

	response.Success(w, responseData, i18n.T(r.Context(), "permission.list_success"))
}
