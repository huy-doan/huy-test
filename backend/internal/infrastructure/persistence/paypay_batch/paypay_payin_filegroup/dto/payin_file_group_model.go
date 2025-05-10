package models

import (
	"time"

	ultis "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch"
)

// PayinFileGroup represents the payin_file_group table
type PayinFileGroup struct {
	ID int `json:"id"`
	ultis.BaseColumnTimestamp

	FileGroupName     string     `json:"file_group_name"`
	PaymentProviderID int        `json:"payment_provider_id"`
	ImportTargetDate  time.Time  `json:"import_target_date"`
	ImportedAt        *time.Time `json:"imported_at"`
}

// TableName specifies the table name for PayinFileGroup
func (PayinFileGroup) TableName() string {
	return "payin_file_group"
}
