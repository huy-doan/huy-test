package outputdata

import model "github.com/huydq/test/internal/domain/model/role"

// BatchUpdateRolePermissionsOutput represents the output for batch permission updates
type BatchUpdateRolePermissionsOutput struct {
	SuccessfulUpdates []int `json:"successful_updates"`
}

// RoleListOutput represents the output for listing roles
type RoleListOutput struct {
	Roles      []*model.Role `json:"roles"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
	Total      int64         `json:"total"`
}
