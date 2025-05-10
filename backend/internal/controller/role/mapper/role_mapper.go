package mapper

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/huydq/test/internal/controller/permission/mapper"
	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	"github.com/labstack/echo/v4"
)

// ExtractIDFromPath extracts the ID parameter from the path
func ExtractIDFromPath(ctx echo.Context) (int, error) {
	idParam := ctx.Param("id")
	if idParam == "" {
		return 0, errors.New("ID parameter is missing")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// MapRequestToCreateRoleInput maps the request to CreateRoleInput
func MapRequestToCreateRoleInput(ctx echo.Context) (*inputdata.CreateRoleInput, error) {
	var req struct {
		Name          string `json:"name"`
		Code          string `json:"code"`
		PermissionIDs []int  `json:"permission_ids"`
	}

	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return nil, err
	}

	// Basic validation
	if req.Name == "" || req.Code == "" {
		return nil, errors.New("name and code are required")
	}

	return &inputdata.CreateRoleInput{
		Name:          req.Name,
		Code:          req.Code,
		PermissionIDs: req.PermissionIDs,
	}, nil
}

// MapRequestToUpdateRoleInput maps the request to UpdateRoleInput
func MapRequestToUpdateRoleInput(ctx echo.Context) (*inputdata.UpdateRoleInput, error) {
	var req struct {
		Name          string `json:"name"`
		Code          string `json:"code"`
		PermissionIDs []int  `json:"permission_ids"`
	}

	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return nil, err
	}

	// Basic validation
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	return &inputdata.UpdateRoleInput{
		Name:          req.Name,
		Code:          req.Code,
		PermissionIDs: req.PermissionIDs,
	}, nil
}

// MapRequestToBatchUpdateRolePermissionsInput maps the request to BatchUpdateRolePermissionsInput
func MapRequestToBatchUpdateRolePermissionsInput(ctx echo.Context) (*inputdata.BatchUpdateRolePermissionsInput, error) {
	var req []struct {
		ID            int   `json:"id"`
		PermissionIDs []int `json:"permission_ids"`
	}

	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return nil, err
	}

	updates := make([]inputdata.RolePermissionUpdate, len(req))
	for i, item := range req {
		updates[i] = inputdata.RolePermissionUpdate{
			ID:            item.ID,
			PermissionIDs: item.PermissionIDs,
		}
	}

	return &inputdata.BatchUpdateRolePermissionsInput{
		Updates: updates, // Changed from RoleUpdates to Updates to match the struct definition
	}, nil
}

// MapRoleListOutputToResponse maps the usecase output to API response format
// Matches the old handler's response format
func MapRoleListOutputToResponse(output *outputdata.RoleListOutput) map[string]interface{} {
	// The old handler returns a map with roles, page, page_size, total_pages, and total
	return map[string]interface{}{
		"roles":       MapRolesToResponse(output.Roles),
		"page":        1,   // Default since we're not paginating in the new implementation
		"page_size":   100, // Default size
		"total_pages": 1,   // Since we return all roles
		"total":       output.Total,
	}
}

// MapRolesToResponse maps a slice of role outputs to response format
func MapRolesToResponse(roles []*outputdata.RoleOutput) []map[string]interface{} {
	if roles == nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, len(roles))
	for i, role := range roles {
		result[i] = MapRoleToResponse(role)
	}
	return result
}

// MapRoleToResponse maps a single role output to response format
func MapRoleToResponse(role *outputdata.RoleOutput) map[string]interface{} {
	if role == nil {
		return nil
	}

	// Create response map that matches serialized role structure in the old handler
	response := map[string]interface{}{
		"id":         role.ID,
		"name":       role.Name,
		"code":       role.Code,
		"created_at": role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), // RFC3339 format
		"updated_at": role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"), // RFC3339 format
	}

	// Include permissions if available
	if role.Permissions != nil {
		// Use the permission mapper to ensure consistent response format
		permissions := make([]map[string]interface{}, 0, len(role.Permissions))
		for _, perm := range role.Permissions {
			permMap := mapper.MapPermissionToResponse(perm)
			permissions = append(permissions, permMap)
		}
		response["permissions"] = permissions
	}

	return response
}

// MapBatchUpdateRolePermissionsOutputToResponse maps the batch update output to response format
func MapBatchUpdateRolePermissionsOutputToResponse(output *outputdata.BatchUpdateRolePermissionsOutput) map[string]interface{} {
	return map[string]interface{}{
		"updated_roles": output.SuccessfulUpdates,
		"total_updated": len(output.SuccessfulUpdates),
	}
}
