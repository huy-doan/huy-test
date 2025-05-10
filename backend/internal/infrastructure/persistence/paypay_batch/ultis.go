package models

import (
	"time"

	"gorm.io/gorm"
)

/** ----------------------------------------------------------
 * Base Column Timestamp
 * ---------------------------------------------------------- */
type BaseColumnTimestamp struct {
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	// swaggerignore: true
	DeletedAt gorm.DeletedAt `json:"deleted_at" swaggerignore:"true"`
}
