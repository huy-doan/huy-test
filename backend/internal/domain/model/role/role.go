package model

import (
	permission "github.com/huydq/test/internal/domain/model/permission"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	object "github.com/huydq/test/internal/domain/object/role"
)

// Role represents a user role in the system
type Role struct {
	ID   int
	Name string
	Code object.RoleCode
	util.BaseColumnTimestamp

	// Relationships
	Permissions []*permission.Permission
}

func (r *Role) IsAdmin() bool {
	return r.Code == object.RoleCodeAdmin
}

func (r *Role) IsNormalUser() bool {
	return r.Code == object.RoleCodeNormalUser
}

func (r *Role) IsBusinessUser() bool {
	return r.Code == object.RoleCodeBusinessUser
}

func (r *Role) IsAccountingUser() bool {
	return r.Code == object.RoleCodeAccoutingUser
}

type NewRoleParams struct {
	ID          int
	Name        string
	Code        object.RoleCode
	Permissions []*permission.Permission
	util.BaseColumnTimestamp
}

func NewRole(params NewRoleParams) *Role {
	return &Role{
		ID:                  params.ID,
		Name:                params.Name,
		Code:                params.Code,
		Permissions:         params.Permissions,
		BaseColumnTimestamp: params.BaseColumnTimestamp,
	}
}
