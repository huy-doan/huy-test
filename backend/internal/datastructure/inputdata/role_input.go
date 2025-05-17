package inputdata

type CreateRoleInput struct {
	Name          string `json:"name" binding:"required"`
	Code          string `json:"code" binding:"required"`
	PermissionIDs []int  `json:"permission_ids"`
}

type UpdateRoleInput struct {
	Name          string `json:"name" binding:"required"`
	PermissionIDs []int  `json:"permission_ids"`
}

type BatchUpdateRolePermissionsInput struct {
	Updates []RolePermissionUpdate `json:"updates" binding:"required"`
}

type RolePermissionUpdate struct {
	ID            int   `json:"id" binding:"required"`
	PermissionIDs []int `json:"permission_ids" binding:"required"`
}
