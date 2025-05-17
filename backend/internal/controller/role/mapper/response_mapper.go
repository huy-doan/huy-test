package mapper

import (
	"github.com/huydq/test/internal/datastructure/outputdata"
	model "github.com/huydq/test/internal/domain/model/role"
	object "github.com/huydq/test/internal/domain/object/basedatetime"
)

// ToRoleResponse maps a single role output to response format
func ToRoleResponse(role *model.Role) *model.RoleResponse {
	if role == nil {
		return nil
	}

	response := &model.RoleResponse{
		Role: model.Role{
			ID:   role.ID,
			Name: role.Name,
			BaseColumnTimestamp: object.BaseColumnTimestamp{
				CreatedAt: role.CreatedAt,
				UpdatedAt: role.UpdatedAt,
			},
		},
	}

	// Include permissions if available
	if role.Permissions != nil {
		response.Permissions = role.Permissions
	}

	return response
}

// ToRoleListResponse maps the usecase output to API response format
func ToRoleListResponse(output *outputdata.RoleListOutput) *model.RoleListResponse {
	if output == nil {
		return &model.RoleListResponse{
			Roles: []model.RoleResponse{},
		}
	}

	roles := make([]model.RoleResponse, 0, len(output.Roles))
	for _, role := range output.Roles {
		if roleResp := ToRoleResponse(role); roleResp != nil {
			roles = append(roles, *roleResp)
		}
	}

	return &model.RoleListResponse{
		Roles: roles,
	}
}

// ToBatchUpdateRolePermissionsResponse maps the batch update output to response format
func ToBatchUpdateRolePermissionsResponse(output *outputdata.BatchUpdateRolePermissionsOutput) *model.BatchUpdateRolePermissionsResponse {
	if output == nil {
		return &model.BatchUpdateRolePermissionsResponse{
			UpdatedRoles: []int{},
			TotalUpdated: 0,
		}
	}

	return &model.BatchUpdateRolePermissionsResponse{
		UpdatedRoles: output.SuccessfulUpdates,
		TotalUpdated: len(output.SuccessfulUpdates),
	}
}
