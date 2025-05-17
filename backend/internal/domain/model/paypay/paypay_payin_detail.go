package model

import (
	"time"

	model "github.com/huydq/test/internal/domain/model/payin"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// PayPayPayinDetail represents the paypay_payin_detail table
type PaypayPayinDetail struct {
	ID int
	util.BaseColumnTimestamp

	PayinFileID          int
	PaymentMerchantID    string
	MerchantBusinessName string
	CutoffDate           *time.Time
	TransactionAmount    float64
	RefundAmount         float64
	UsageFee             float64
	PlatformFee          float64
	InitialFee           float64
	Tax                  float64
	Cashback             float64
	Adjustment           float64
	Fee                  float64
	Amount               float64

	PayinFile model.PayinFile
}
