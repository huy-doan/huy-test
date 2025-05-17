package model

import (
	"time"

	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// PayinFileGroup represents the payin_file_group table
type PayinFileGroup struct {
	ID int `json:"id"`
	util.BaseColumnTimestamp

	FileGroupName     string
	PaymentProviderID int
	ImportTargetDate  time.Time
	ImportedAt        *time.Time
}
