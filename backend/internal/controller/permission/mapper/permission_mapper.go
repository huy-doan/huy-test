package mapper

import (
	"strconv"
	"strings"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	"github.com/labstack/echo/v4"
)

// MapRequestToGetPermissionsByIDsInput converts Echo request to GetPermissionsByIDsInput
func MapRequestToGetPermissionsByIDsInput(ctx echo.Context) (*inputdata.GetPermissionsByIDsInput, error) {
	idsParam := ctx.QueryParam("ids")
	if idsParam == "" {
		return &inputdata.GetPermissionsByIDsInput{IDs: []int{}}, nil
	}

	idStrs := strings.Split(idsParam, ",")
	ids := make([]int, 0, len(idStrs))

	for _, idStr := range idStrs {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return &inputdata.GetPermissionsByIDsInput{IDs: ids}, nil
}

// MapPermissionListOutputToResponse converts ListPermissionOutput to response format
// Matches the old handler's response format
func MapPermissionListOutputToResponse(output *outputdata.PermissionListOutput) map[string]interface{} {
	// The old handler returns a map with a single "permissions" key
	return map[string]interface{}{
		"permissions": MapPermissionsToResponse(output.Permissions),
	}
}

// MapPermissionsToResponse converts slice of PermissionOutput to response format
func MapPermissionsToResponse(permissions []*outputdata.PermissionOutput) []map[string]interface{} {
	if permissions == nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, len(permissions))
	for i, permission := range permissions {
		result[i] = MapPermissionToResponse(permission)
	}
	return result
}

// MapPermissionToResponse converts a single PermissionOutput to response format
func MapPermissionToResponse(permission *outputdata.PermissionOutput) map[string]interface{} {
	if permission == nil {
		return nil
	}

	// Create response map that matches serialized permission structure in the old handler
	response := map[string]interface{}{
		"id":         permission.ID,
		"name":       permission.Name,
		"code":       permission.Code,
		"screen_id":  permission.ScreenID,
		"created_at": permission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), // RFC3339 format
		"updated_at": permission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"), // RFC3339 format
	}

	// Include screen if available
	if permission.Screen != nil {
		response["screen"] = map[string]interface{}{
			"id":          permission.Screen.ID,
			"name":        permission.Screen.Name,
			"screen_code": permission.Screen.ScreenCode,
			"screen_path": permission.Screen.ScreenPath,
			"created_at":  permission.Screen.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			"updated_at":  permission.Screen.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return response
}
