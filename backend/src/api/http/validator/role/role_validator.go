package role

import (
	"github.com/vnlab/makeshop-payment/src/lib/validator"
)

// CreateRoleRequest represents a request to create a new role
type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Code          string `json:"code" binding:"required"`
	PermissionIDs []int  `json:"permission_ids"`
}

// UpdateRoleRequest represents a request to update an existing role
type UpdateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	PermissionIDs []int  `json:"permission_ids"`
}

// RolePermissionUpdateItem represents a single role's permissions update
type RolePermissionUpdateItem struct {
	ID            int   `json:"id" binding:"required"`
	PermissionIDs []int `json:"permission_ids" binding:"required"`
}

// BatchUpdateRolePermissionsRequest represents a request to update permissions for multiple roles
type BatchUpdateRolePermissionsRequest []RolePermissionUpdateItem

// Validate validates the create role request
func (r *CreateRoleRequest) Validate() error {
	validate := validator.GetValidate()
	return validate.Struct(r)
}

// Validate validates the update role request
func (r *UpdateRoleRequest) Validate() error {
	validate := validator.GetValidate()
	return validate.Struct(r)
}

// Validate validates the batch update role permissions request
func (r *BatchUpdateRolePermissionsRequest) Validate() error {
	validate := validator.GetValidate()
	return validate.Var(r, "required,dive")
}
