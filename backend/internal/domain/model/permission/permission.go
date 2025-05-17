package model

import (
	screen "github.com/huydq/test/internal/domain/model/screen"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	object "github.com/huydq/test/internal/domain/object/permission"
)

// Permission represents a permission in the system
type Permission struct {
	ID       int                   `json:"id"`
	Name     string                `json:"name" binding:"required"`
	Code     object.PermissionCode `json:"code" binding:"required"`
	ScreenID int                   `json:"screen_id" binding:"required"`
	Screen   *screen.Screen        `json:"screen" binding:"required"`
	util.BaseColumnTimestamp
}

type NewPermissionParams struct {
	ID int
	util.BaseColumnTimestamp
	Name     string
	Code     object.PermissionCode
	ScreenID int
	Screen   *screen.Screen
}

func NewPermission(params NewPermissionParams) *Permission {
	return &Permission{
		Name:                params.Name,
		Code:                params.Code,
		ScreenID:            params.ScreenID,
		Screen:              params.Screen,
		BaseColumnTimestamp: params.BaseColumnTimestamp,
	}
}

type PermissionListResponse struct {
	Permissions []*Permission `json:"permissions"`
	Total       int64         `json:"total"`
}
