package model

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// Screen represents a system screen
type Screen struct {
	ID         int    `json:"id"`
	Name       string `json:"name" binding:"required"`
	ScreenCode string `json:"screen_code" binding:"required"`
	ScreenPath string `json:"screen_path" binding:"required"`
	util.BaseColumnTimestamp
}

type NewScreenParams struct {
	ID int
	util.BaseColumnTimestamp
	Name       string
	ScreenCode string
	ScreenPath string
}

func NewScreen(params NewScreenParams) *Screen {
	return &Screen{
		ID:                  params.ID,
		Name:                params.Name,
		ScreenCode:          params.ScreenCode,
		ScreenPath:          params.ScreenPath,
		BaseColumnTimestamp: params.BaseColumnTimestamp,
	}
}
