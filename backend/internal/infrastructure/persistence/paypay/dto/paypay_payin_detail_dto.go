package dto

import (
	"time"

	model "github.com/huydq/test/internal/domain/model/paypay"
	"github.com/huydq/test/internal/infrastructure/persistence/payin/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

// PayPayPayinDetail represents the paypay_payin_detail table
type PaypayPayinDetail struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`
	persistence.BaseColumnTimestamp

	PayinFileID          int        `json:"payin_file_id"`
	PaymentMerchantID    string     `json:"payment_merchant_id"`
	MerchantBusinessName string     `json:"merchant_business_name"`
	CutoffDate           *time.Time `json:"cutoff_date"`
	TransactionAmount    float64    `json:"transaction_amount"`
	RefundAmount         float64    `json:"refund_amount"`
	UsageFee             float64    `json:"usage_fee"`
	PlatformFee          float64    `json:"platform_fee"`
	InitialFee           float64    `json:"initial_fee"`
	Tax                  float64    `json:"tax"`
	Cashback             float64    `json:"cashback"`
	Adjustment           float64    `json:"adjustment"`
	Fee                  float64    `json:"fee"`
	Amount               float64    `json:"amount"`

	PayinFile *dto.PayinFile `json:"payin_file,omitempty"`
}

// TableName specifies the table name for PayPayPayinDetail
func (PaypayPayinDetail) TableName() string {
	return "paypay_payin_detail"
}

func (dto *PaypayPayinDetail) ToPaypayPayinDetailModel() *model.PaypayPayinDetail {
	paypayPayinDetailModel := &model.PaypayPayinDetail{
		ID:                   dto.ID,
		PayinFileID:          dto.PayinFileID,
		PaymentMerchantID:    dto.PaymentMerchantID,
		MerchantBusinessName: dto.MerchantBusinessName,
		CutoffDate:           dto.CutoffDate,
		TransactionAmount:    dto.TransactionAmount,
		RefundAmount:         dto.RefundAmount,
		UsageFee:             dto.UsageFee,
		PlatformFee:          dto.PlatformFee,
		InitialFee:           dto.InitialFee,
		Tax:                  dto.Tax,
		Cashback:             dto.Cashback,
		Adjustment:           dto.Adjustment,
		Fee:                  dto.Fee,
	}
	paypayPayinDetailModel.CreatedAt = dto.CreatedAt
	paypayPayinDetailModel.UpdatedAt = dto.UpdatedAt
	return paypayPayinDetailModel
}

func ToPaypayPayinDetailDTO(p *model.PaypayPayinDetail) *PaypayPayinDetail {
	paypayPayinDetailDTO := &PaypayPayinDetail{
		ID:                   p.ID,
		PayinFileID:          p.PayinFileID,
		PaymentMerchantID:    p.PaymentMerchantID,
		MerchantBusinessName: p.MerchantBusinessName,
		CutoffDate:           p.CutoffDate,
		TransactionAmount:    p.TransactionAmount,
		RefundAmount:         p.RefundAmount,
		UsageFee:             p.UsageFee,
		PlatformFee:          p.PlatformFee,
		InitialFee:           p.InitialFee,
		Tax:                  p.Tax,
		Cashback:             p.Cashback,
		Adjustment:           p.Adjustment,
		Fee:                  p.Fee,
	}
	paypayPayinDetailDTO.CreatedAt = p.CreatedAt
	paypayPayinDetailDTO.UpdatedAt = p.UpdatedAt
	return paypayPayinDetailDTO
}

func ToPaypayPayinDetailDTOs(paypayPayinDetails []*model.PaypayPayinDetail) []*PaypayPayinDetail {
	paypayPayinDetailDTOs := make([]*PaypayPayinDetail, len(paypayPayinDetails))
	for i, paypayPayinDetail := range paypayPayinDetails {
		paypayPayinDetailDTOs[i] = ToPaypayPayinDetailDTO(paypayPayinDetail)
	}
	return paypayPayinDetailDTOs
}

func ToPaypayPayinDetailModels(paypayPayinDetailDTOs []*PaypayPayinDetail) []*model.PaypayPayinDetail {
	paypayPayinDetails := make([]*model.PaypayPayinDetail, len(paypayPayinDetailDTOs))
	for i, paypayPayinDetailDTO := range paypayPayinDetailDTOs {
		paypayPayinDetails[i] = paypayPayinDetailDTO.ToPaypayPayinDetailModel()
	}
	return paypayPayinDetails
}
