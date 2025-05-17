package models

import (
	"time"
)

// PayPayPayinSummary represents the paypay_payin_summary table
type PayPayPayinSummary struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	PayinFileID       int        `json:"payin_file_id"`
	CorporateName     string     `json:"corporate_name"`
	CutoffDate        *time.Time `json:"cutoff_date"`
	PaymentDate       *time.Time `json:"payment_date"`
	TransactionAmount float64    `json:"transaction_amount"`
	RefundAmount      float64    `json:"refund_amount"`
	UsageFee          float64    `json:"usage_fee"`
	PlatformFee       float64    `json:"platform_fee"`
	InitialFee        float64    `json:"initial_fee"`
	Tax               float64    `json:"tax"`
	Cashback          float64    `json:"cashback"`
	Adjustment        float64    `json:"adjustment"`
	Fee               float64    `json:"fee"`
	Amount            float64    `json:"amount"`

	PayinFile *PayinFile `json:"payin_file,omitempty"`
}

// TableName specifies the table name for PayPayPayinSummary
func (PayPayPayinSummary) TableName() string {
	return "paypay_payin_summary"
}
