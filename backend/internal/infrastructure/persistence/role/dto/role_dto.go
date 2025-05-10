package dto

import (
	permissionDto "github.com/huydq/test/internal/infrastructure/persistence/permission/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

type RoleDTO struct {
	ID int `json:"id"`
	persistence.BaseColumnTimestamp

	Name string `json:"name"`
	Code string `json:"code"`

	Permissions []*permissionDto.PermissionDTO `json:"permissions" gorm:"many2many:role_permission;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:PermissionID"`
}

// TableName specifies the database table name
func (RoleDTO) TableName() string {
	return "role"
}
