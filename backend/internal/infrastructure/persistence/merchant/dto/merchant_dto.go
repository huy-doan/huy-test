package dto

import (
	"github.com/huydq/test/internal/domain/model/merchant"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

// MerchantDTO is the database representation of a merchant
type MerchantDTO struct {
	ID                int    `gorm:"column:id;primaryKey"`
	PaymentProviderID int    `gorm:"column:payment_provider_id"`
	PaymentMerchantID string `gorm:"column:payment_merchant_id"`
	MerchantName      string `gorm:"column:merchant_name"`
	ShopID            int    `gorm:"column:shop_id"`
	ShopURL           string `gorm:"column:shop_url"`

	persistence.BaseColumnTimestamp

	// Relations
	MerchantPaymentProviderReview *PaymentProviderReviewDTO `gorm:"foreignKey:MerchantID"`
}

// TableName specifies the table name for MerchantDTO
func (MerchantDTO) TableName() string {
	return "merchant"
}

// PaymentProviderReviewDTO is the database representation of a payment provider review
type PaymentProviderReviewDTO struct {
	ID                   int `gorm:"column:id;primaryKey"`
	MerchantID           int `gorm:"column:merchant_id"`
	MerchantReviewStatus int `gorm:"column:merchant_review_status"`

	persistence.BaseColumnTimestamp
}

// TableName specifies the table name for PaymentProviderReviewDTO
func (PaymentProviderReviewDTO) TableName() string {
	return "payment_provider_review"
}

// ToMerchantModel converts a MerchantDTO to a Merchant model
func (dto *MerchantDTO) ToMerchantModel() *merchant.Merchant {
	result := &merchant.Merchant{
		ID:                dto.ID,
		PaymentProviderID: dto.PaymentProviderID,
		PaymentMerchantID: dto.PaymentMerchantID,
		MerchantName:      dto.MerchantName,
		ShopID:            dto.ShopID,
		ShopURL:           dto.ShopURL,
	}

	if dto.MerchantPaymentProviderReview != nil {
		result.MerchantPaymentProviderReview = &merchant.PaymentProviderReview{
			ID:                   dto.MerchantPaymentProviderReview.ID,
			MerchantID:           dto.MerchantPaymentProviderReview.MerchantID,
			MerchantReviewStatus: dto.MerchantPaymentProviderReview.MerchantReviewStatus,
		}
	}

	return result
}

// ToMerchantDTO converts a Merchant model to a MerchantDTO
func ToMerchantDTO(m *merchant.Merchant) *MerchantDTO {
	result := &MerchantDTO{
		ID:                m.ID,
		PaymentProviderID: m.PaymentProviderID,
		PaymentMerchantID: m.PaymentMerchantID,
		MerchantName:      m.MerchantName,
		ShopID:            m.ShopID,
		ShopURL:           m.ShopURL,
	}

	if m.MerchantPaymentProviderReview != nil {
		result.MerchantPaymentProviderReview = &PaymentProviderReviewDTO{
			ID:                   m.MerchantPaymentProviderReview.ID,
			MerchantID:           m.MerchantPaymentProviderReview.MerchantID,
			MerchantReviewStatus: m.MerchantPaymentProviderReview.MerchantReviewStatus,
		}
	}

	return result
}
