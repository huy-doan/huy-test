package dto

import persistence "github.com/huydq/test/internal/infrastructure/persistence/util"

type Screen struct {
	ID int `json:"id"`
	persistence.BaseColumnTimestamp

	Name       string `json:"name"`
	ScreenCode string `json:"screen_code"`
	ScreenPath string `json:"screen_path"`
}

func (Screen) TableName() string {
	return "screen"
}
