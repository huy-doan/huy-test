package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/huydq/test/src/api/http/errors"
	"github.com/huydq/test/src/api/http/middleware"
	"github.com/huydq/test/src/api/http/response"
	"github.com/huydq/test/src/api/http/serializers"
	validator "github.com/huydq/test/src/api/http/validator/role"
	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories/filter"
	"github.com/huydq/test/src/infrastructure/logger"
	"github.com/huydq/test/src/lib/i18n"
	"github.com/huydq/test/src/lib/utils"
	"github.com/huydq/test/src/usecase"
)

type RoleHandler struct {
	roleUsecase *usecase.RoleUsecase
}

func NewRoleHandler(roleUsecase *usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{
		roleUsecase: roleUsecase,
	}
}

// CreateRole handles creating a new role
// @Summary Create a new role
// @Description Create a new role with permissions
// @Tags Admin Role Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body role.CreateRoleRequest true "Role creation details"
// @Success 201 {object} response.Response "Created"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/roles [post]
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req validator.CreateRoleRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Code == "" {
		response.BadRequest(w, i18n.T(r.Context(), "validation.required_fields"), nil)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	// Create role model
	role := &models.Role{
		Name: req.Name,
		Code: req.Code,
	}

	// Add permissions if provided
	if len(req.PermissionIDs) > 0 {
		permissions := make([]*models.Permission, 0, len(req.PermissionIDs))
		for _, id := range req.PermissionIDs {
			permissions = append(permissions, &models.Permission{ID: id})
		}
		role.Permissions = permissions
	}

	if err := h.roleUsecase.CreateRole(r.Context(), role); err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Created(w, serializers.NewRoleSerializer(role).Serialize(), i18n.T(r.Context(), "role.create_success"))
}

// GetRoleByID handles getting a role by ID
// @Summary Get a role by ID
// @Description Get role details by ID
// @Tags Admin Role Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} response.Response{data=models.Role}
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Not Found"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/roles/{id} [get]
func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.BadRequest(w, i18n.T(r.Context(), "params.not_found", "ID"), nil)
		return
	}
	idStr := pathParts[len(pathParts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, i18n.T(r.Context(), "params.invalid_number", "ID"), nil)
		return
	}

	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	role, err := h.roleUsecase.GetRoleByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if role == nil {
		response.NotFound(w, i18n.T(r.Context(), "role.not_found"))
		return
	}

	response.Success(w, serializers.NewRoleSerializer(role).Serialize(), i18n.T(r.Context(), "role.get_success"))
}

// UpdateRole handles updating a role
// @Summary Update a role
// @Description Update a role's details
// @Tags Admin Role Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Param request body role.UpdateRoleRequest true "Role update details"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Not Found"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/roles/{id} [put]
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.BadRequest(w, i18n.T(r.Context(), "params.not_found", "ID"), nil)
		return
	}
	idStr := pathParts[len(pathParts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, i18n.T(r.Context(), "params.invalid_number", "ID"), nil)
		return
	}

	var req validator.UpdateRoleRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	// Check if the role exists first
	existingRole, err := h.roleUsecase.GetRoleByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if existingRole == nil {
		response.NotFound(w, i18n.T(r.Context(), "role.not_found"))
		return
	}

	// Convert request to usecase format
	updateData := &usecase.UpdateRoleRequest{
		Name:          req.Name,
		PermissionIDs: req.PermissionIDs,
	}

	if err := h.roleUsecase.UpdateRole(r.Context(), id, updateData); err != nil {
		h.handleError(w, r, err)
		return
	}

	// Get the updated role for the response
	updatedRole, err := h.roleUsecase.GetRoleByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Success(w, serializers.NewRoleSerializer(updatedRole).Serialize(), i18n.T(r.Context(), "role.update_success"))
}

// DeleteRole handles deleting a role
// @Summary Delete a role
// @Description Delete a role by ID
// @Tags Admin Role Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Not Found"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		response.BadRequest(w, i18n.T(r.Context(), "params.not_found", "ID"), nil)
		return
	}
	idStr := pathParts[len(pathParts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(w, i18n.T(r.Context(), "params.invalid_number", "ID"), nil)
		return
	}

	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	// Check if the role exists
	role, err := h.roleUsecase.GetRoleByID(r.Context(), id)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	if role == nil {
		response.NotFound(w, i18n.T(r.Context(), "role.not_found"))
		return
	}

	if err := h.roleUsecase.DeleteRole(r.Context(), id); err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Success(w, nil, i18n.T(r.Context(), "role.delete_success"))
}

// ListRoles handles listing roles with pagination
// @Summary List roles
// @Description Get a list of roles with pagination
// @Tags Admin Role Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "Page number (default: 1)"
// @Param page_size query integer false "Page size (default: 10, max: 100)"
// @Success 200 {object} response.Response "Success"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/roles [get]
func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	roleFilter := filter.NewRoleFilter()
	baseFilter := utils.ExtractPaginationAndSorting(r)
	roleFilter.Pagination = baseFilter.Pagination
	roleFilter.Sort = baseFilter.Sort

	page, pageSize := utils.ExtractPaginationParams(r)
	err := applyRoleFilterParams(r, roleFilter)

	if err != nil {
		response.ValidationError(w, err)
		return
	}

	roleFilter.ApplyFilters()

	roles, totalPages, total, err := h.roleUsecase.ListRoles(r.Context(), roleFilter)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	responseData := map[string]any{
		"roles":       serializers.SerializeRoleCollection(roles),
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
		"total":       total,
	}

	response.Success(w, responseData, i18n.T(r.Context(), "role.list_success"))
}

func applyRoleFilterParams(r *http.Request, roleFilter *filter.RoleFilter) error {
	name := utils.ExtractStringParam(r, "name")
	roleFilter.Name = name

	return nil
}

// BatchUpdateRolePermissions handles updating permissions for multiple roles at once
// @Summary Batch update role permissions
// @Description Update permissions for multiple roles in a single request
// @Tags Admin Role Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body role.BatchUpdateRolePermissionsRequest true "Batch role permissions update data"
// @Success 200 {object} response.Response "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /admin/roles/permissions/batch [post]
func (h *RoleHandler) BatchUpdateRolePermissions(w http.ResponseWriter, r *http.Request) {
	var req validator.BatchUpdateRolePermissionsRequest
	if err := utils.ParseJSONBody(r, &req); err != nil {
		response.ValidationError(w, err)
		return
	}

	if err := req.Validate(); err != nil {
		response.ValidationError(w, err)
		return
	}

	// Check admin role
	roleCode, ok := r.Context().Value(middleware.RoleCodeKey).(string)
	if !ok || roleCode != string(models.RoleCodeAdmin) {
		response.Forbidden(w, i18n.T(r.Context(), "account.unauthorized"))
		return
	}

	// Convert request to usecase format
	updates := make([]struct {
		ID            int
		PermissionIDs []int
	}, len(req))

	for i, item := range req {
		updates[i] = struct {
			ID            int
			PermissionIDs []int
		}{
			ID:            item.ID,
			PermissionIDs: item.PermissionIDs,
		}
	}

	// Process the batch update
	updatedIDs, err := h.roleUsecase.BatchUpdateRolePermissions(r.Context(), updates)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Success(w, map[string]any{
		"updated_roles": updatedIDs,
		"total_updated": len(updatedIDs),
	}, i18n.T(r.Context(), "role.batch_update_success"))
}

// handleError centralizes error handling for role-related operations
func (h *RoleHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	switch err.Error() {
	case "one or more permission IDs do not exist":
		response.BadRequest(w, i18n.T(r.Context(), "permission.not_found"), nil)
	case "role not found":
		response.NotFound(w, i18n.T(r.Context(), "role.not_found"))
	case "role name already exists":
		response.BadRequest(w, i18n.T(r.Context(), "role.duplicate_name"), nil)
	case "role code already exists":
		response.BadRequest(w, i18n.T(r.Context(), "role.duplicate_code"), nil)
	default:
		logger.GetLogger().Error(fmt.Sprintf("Role handler error: %v", err), nil)
		response.Error(w, errors.InternalError(i18n.T(r.Context(), "common.error")))
	}
}
