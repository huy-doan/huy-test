package dto

import (
	screenDto "github.com/huydq/test/internal/infrastructure/persistence/screen/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

// Role represents a user role in the system
type Permission struct {
	ID int `json:"id"`
	persistence.BaseColumnTimestamp

	Name     string            `json:"name"`
	Code     string            `json:"code"`
	ScreenID int               `json:"screen_id"`
	Screen   *screenDto.Screen `json:"screen" gorm:"foreignKey:ScreenID"`
}

// TableName specifies the database table name
func (Permission) TableName() string {
	return "permission"
}
