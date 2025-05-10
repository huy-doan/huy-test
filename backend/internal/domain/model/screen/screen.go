package model

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// Screen represents a system screen
type Screen struct {
	ID         int
	Name       string
	ScreenCode string
	ScreenPath string
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
