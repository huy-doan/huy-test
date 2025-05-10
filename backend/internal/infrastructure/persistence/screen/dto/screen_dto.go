package dto

import persistence "github.com/huydq/test/internal/infrastructure/persistence/util"

type ScreenDTO struct {
	ID int `json:"id"`
	persistence.BaseColumnTimestamp

	Name       string `json:"name"`
	ScreenCode string `json:"screen_code"`
	ScreenPath string `json:"screen_path"`
}

func (ScreenDTO) TableName() string {
	return "screen"
}
