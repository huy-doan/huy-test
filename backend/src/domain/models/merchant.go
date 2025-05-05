package models

// Merchant represents the merchant table
type Merchant struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	PaymentProviderID int    `json:"payment_provider_id"`
	PaymentMerchantID string `json:"payment_merchant_id"`
	MerchantName      string `json:"merchant_name"`
	ShopID            int    `json:"shop_id"`
	ShopURL           string `json:"shop_url"`

	// Relations
	MerchantPaymentProviderReview *PaymentProviderReview `json:"merchant_payment_provider_review" gorm:"foreignKey:MerchantID"`
}

// TableName specifies the table name for Merchant
func (Merchant) TableName() string {
	return "merchant"
}
