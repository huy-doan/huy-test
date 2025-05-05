package models

import "time"

// PaymentProviderReview represents the payment_provider_review table
type PaymentProviderReview struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	PaymentProviderID    int       `json:"payment_provider_id"`
	MerchantReviewStatus int       `json:"merchant_review_status"`
	PaymentMerchantID    *int      `json:"payment_merchant_id"`
	MerchantID           *int      `json:"merchant_id" gorm:"column:merchant_id"`
	RegisteredAt         time.Time `json:"registered_at"`

	Merchant *Merchant `json:"merchant" gorm:"foreignKey:MerchantID"`
}

const MerchantReviewStatusUnderReview = 1
const MerchantReviewStatusApproved = 2
const MerchantReviewStatusDenied = 3

// TableName specifies the table name for PaymentProviderReview
func (PaymentProviderReview) TableName() string {
	return "payment_provider_review"
}
