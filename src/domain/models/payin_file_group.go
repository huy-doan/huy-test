package models

import (
	"time"
)

// PayinFileGroup represents the payin_file_group table
type PayinFileGroup struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	PaymentProviderID int        `json:"payment_provider_id"`
	ImportTargetDate  time.Time  `json:"import_target_date"`
	ImportedAt        *time.Time `json:"imported_at"`
}

// TableName specifies the table name for PayinFileGroup
func (PayinFileGroup) TableName() string {
	return "payin_file_group"
}
