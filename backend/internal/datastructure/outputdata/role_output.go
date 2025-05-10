package outputdata

import (
	"time"
)

// RoleOutput represents the output data for a role
type RoleOutput struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Code        string              `json:"code"`
	Permissions []*PermissionOutput `json:"permissions,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// BatchUpdateRolePermissionsOutput represents the output for batch permission updates
type BatchUpdateRolePermissionsOutput struct {
	SuccessfulUpdates []int `json:"successful_updates"`
}

// RoleListOutput represents the output for listing roles
type RoleListOutput struct {
	Roles      []*RoleOutput `json:"roles"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
	Total      int64         `json:"total"`
}
