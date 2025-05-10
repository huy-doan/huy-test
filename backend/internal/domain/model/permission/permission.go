package model

import (
	screen "github.com/huydq/test/internal/domain/model/screen"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	object "github.com/huydq/test/internal/domain/object/permission"
)

// Permission represents a permission in the system
type Permission struct {
	ID       int
	Name     string
	Code     object.PermissionCode
	ScreenID int
	Screen   *screen.Screen
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
