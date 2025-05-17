package dto

import (
	model "github.com/huydq/test/internal/domain/model/merchant"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

// Merchant is the database representation of a merchant
type Merchant struct {
	ID                int    `gorm:"column:id;primaryKey"`
	PaymentProviderID int    `gorm:"column:payment_provider_id"`
	PaymentMerchantID string `gorm:"column:payment_merchant_id"`
	MerchantName      string `gorm:"column:merchant_name"`
	ShopID            int    `gorm:"column:shop_id"`
	ShopURL           string `gorm:"column:shop_url"`

	persistence.BaseColumnTimestamp

	// Relations
	MerchantPaymentProviderReview *PaymentProviderReview `gorm:"foreignKey:MerchantID"`
}

// TableName specifies the table name for Merchant
func (Merchant) TableName() string {
	return "merchant"
}

// PaymentProviderReview is the database representation of a payment provider review
type PaymentProviderReview struct {
	ID                   int `gorm:"column:id;primaryKey"`
	MerchantID           int `gorm:"column:merchant_id"`
	MerchantReviewStatus int `gorm:"column:merchant_review_status"`

	persistence.BaseColumnTimestamp
}

// TableName specifies the table name for PaymentProviderReview
func (PaymentProviderReview) TableName() string {
	return "payment_provider_review"
}

// ToMerchantModel converts a Merchant to a Merchant model
func (dto *Merchant) ToMerchantModel() *model.Merchant {
	result := &model.Merchant{
		ID:                dto.ID,
		PaymentProviderID: dto.PaymentProviderID,
		PaymentMerchantID: dto.PaymentMerchantID,
		MerchantName:      dto.MerchantName,
		ShopID:            dto.ShopID,
		ShopURL:           dto.ShopURL,
	}

	result.CreatedAt = dto.CreatedAt
	result.UpdatedAt = dto.UpdatedAt

	if dto.MerchantPaymentProviderReview != nil {
		result.MerchantPaymentProviderReview = &model.PaymentProviderReview{
			ID:                   dto.MerchantPaymentProviderReview.ID,
			MerchantID:           dto.MerchantPaymentProviderReview.MerchantID,
			MerchantReviewStatus: dto.MerchantPaymentProviderReview.MerchantReviewStatus,
		}
	}

	return result
}

// ToMerchantDTO converts a Merchant model to a Merchant
func ToMerchantDTO(m *model.Merchant) *Merchant {
	result := &Merchant{
		ID:                m.ID,
		PaymentProviderID: m.PaymentProviderID,
		PaymentMerchantID: m.PaymentMerchantID,
		MerchantName:      m.MerchantName,
		ShopID:            m.ShopID,
		ShopURL:           m.ShopURL,
	}

	result.CreatedAt = m.CreatedAt
	result.UpdatedAt = m.UpdatedAt

	if m.MerchantPaymentProviderReview != nil {
		result.MerchantPaymentProviderReview = &PaymentProviderReview{
			ID:                   m.MerchantPaymentProviderReview.ID,
			MerchantID:           m.MerchantPaymentProviderReview.MerchantID,
			MerchantReviewStatus: m.MerchantPaymentProviderReview.MerchantReviewStatus,
		}
	}

	return result
}
