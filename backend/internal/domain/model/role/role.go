package model

import (
	permission "github.com/huydq/test/internal/domain/model/permission"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	permissionObject "github.com/huydq/test/internal/domain/object/permission"
)

// Role represents a user role in the system
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	util.BaseColumnTimestamp

	// Relationships
	Permissions []*permission.Permission `json:"permissions"`
}

// HasPermission checks if the role has the specified permission
func (r *Role) HasPermission(permissions ...permissionObject.PermissionCode) bool {
	if r.Permissions == nil || len(r.Permissions) == 0 {
		return false
	}

	// Convert permissions to map for faster lookup
	requiredPerms := make(map[permissionObject.PermissionCode]bool)
	for _, p := range permissions {
		requiredPerms[p] = true
	}

	// Check if the role has any of the required permissions
	for _, perm := range r.Permissions {
		if requiredPerms[perm.Code] {
			return true
		}
	}

	return false
}

type NewRoleParams struct {
	ID          int
	Name        string
	Permissions []*permission.Permission
	util.BaseColumnTimestamp
}

func NewRole(params NewRoleParams) *Role {
	return &Role{
		ID:                  params.ID,
		Name:                params.Name,
		Permissions:         params.Permissions,
		BaseColumnTimestamp: params.BaseColumnTimestamp,
	}
}

type RolePermissionUpdateItem struct {
	ID            int   `json:"id" binding:"required"`
	PermissionIDs []int `json:"permission_ids" binding:"required"`
}

type BatchUpdateRolePermissionsRequest []RolePermissionUpdateItem

type RoleResponse struct {
	Role
}

// RoleListResponse represents the API response for a list of roles
type RoleListResponse struct {
	Roles []RoleResponse `json:"roles"`
}

// BatchUpdateRolePermissionsResponse represents the API response for batch updating role permissions
type BatchUpdateRolePermissionsResponse struct {
	UpdatedRoles []int `json:"updated_roles"`
	TotalUpdated int   `json:"total_updated"`
}
