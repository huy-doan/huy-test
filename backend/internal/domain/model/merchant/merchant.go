package model

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// Merchant represents the merchant entity
type Merchant struct {
	ID                int    `json:"id"`
	PaymentProviderID int    `json:"payment_provider_id"`
	PaymentMerchantID string `json:"payment_merchant_id"`
	MerchantName      string `json:"merchant_name"`
	ShopID            int    `json:"shop_id"`
	ShopURL           string `json:"shop_url"`

	util.BaseColumnTimestamp

	MerchantPaymentProviderReview *PaymentProviderReview `json:"merchant_payment_provider_review"`
}

// PaymentProviderReview represents the payment provider review entity
type PaymentProviderReview struct {
	ID                   int `json:"id"`
	MerchantID           int `json:"merchant_id"`
	MerchantReviewStatus int `json:"merchant_review_status"`

	util.BaseColumnTimestamp
}
