package inputdata

// CreateRoleInput represents the input data for creating a new role
type CreateRoleInput struct {
	Name          string `json:"name" binding:"required"`
	Code          string `json:"code" binding:"required"`
	PermissionIDs []int  `json:"permission_ids"`
}

// UpdateRoleInput represents the input data for updating a role
type UpdateRoleInput struct {
	Name          string `json:"name" binding:"required"`
	Code          string `json:"code" binding:"required"`
	PermissionIDs []int  `json:"permission_ids"`
}

// BatchUpdateRolePermissionsInput represents the input for updating permissions for multiple roles
type BatchUpdateRolePermissionsInput struct {
	Updates []RolePermissionUpdate `json:"updates" binding:"required"`
}

// RolePermissionUpdate represents a single role permission update in a batch
type RolePermissionUpdate struct {
	ID            int   `json:"id" binding:"required"`
	PermissionIDs []int `json:"permission_ids" binding:"required"`
}
