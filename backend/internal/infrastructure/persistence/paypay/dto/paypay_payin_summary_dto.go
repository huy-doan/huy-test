package dto

import (
	"time"

	model "github.com/huydq/test/internal/domain/model/paypay"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"

	"github.com/huydq/test/internal/infrastructure/persistence/payin/dto"
)

// PayPayPayinSummary represents the paypay_payin_summary table
type PaypayPayinSummary struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`
	persistence.BaseColumnTimestamp

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

	PayinFile *dto.PayinFile `json:"payin_file,omitempty"`
}

// TableName specifies the table name for PayPayPayinSummary
func (PaypayPayinSummary) TableName() string {
	return "paypay_payin_summary"
}

func (dto *PaypayPayinSummary) ToPaypayPayinSummaryModel() *model.PaypayPayinSummary {
	paypayPayinSummaryModel := &model.PaypayPayinSummary{
		ID:                dto.ID,
		PayinFileID:       dto.PayinFileID,
		CorporateName:     dto.CorporateName,
		CutoffDate:        dto.CutoffDate,
		PaymentDate:       dto.PaymentDate,
		TransactionAmount: dto.TransactionAmount,
		RefundAmount:      dto.RefundAmount,
		UsageFee:          dto.UsageFee,
		PlatformFee:       dto.PlatformFee,
		InitialFee:        dto.InitialFee,
		Tax:               dto.Tax,
		Cashback:          dto.Cashback,
		Adjustment:        dto.Adjustment,
		Fee:               dto.Fee,
	}
	paypayPayinSummaryModel.CreatedAt = dto.CreatedAt
	paypayPayinSummaryModel.UpdatedAt = dto.UpdatedAt
	return paypayPayinSummaryModel
}

func ToPaypayPayinSummaryDTO(paypayPayinSummaries *model.PaypayPayinSummary) *PaypayPayinSummary {
	paypayPayinSummaryDTO := &PaypayPayinSummary{
		ID:                paypayPayinSummaries.ID,
		PayinFileID:       paypayPayinSummaries.PayinFileID,
		CorporateName:     paypayPayinSummaries.CorporateName,
		CutoffDate:        paypayPayinSummaries.CutoffDate,
		PaymentDate:       paypayPayinSummaries.PaymentDate,
		TransactionAmount: paypayPayinSummaries.TransactionAmount,
		RefundAmount:      paypayPayinSummaries.RefundAmount,
		UsageFee:          paypayPayinSummaries.UsageFee,
		PlatformFee:       paypayPayinSummaries.PlatformFee,
		InitialFee:        paypayPayinSummaries.InitialFee,
		Tax:               paypayPayinSummaries.Tax,
		Cashback:          paypayPayinSummaries.Cashback,
		Adjustment:        paypayPayinSummaries.Adjustment,
		Fee:               paypayPayinSummaries.Fee,
	}
	paypayPayinSummaryDTO.CreatedAt = paypayPayinSummaries.CreatedAt
	paypayPayinSummaryDTO.UpdatedAt = paypayPayinSummaries.UpdatedAt
	return paypayPayinSummaryDTO
}

func ToPaypayPayinSummaryDTOs(paypayPayinSummaries []*model.PaypayPayinSummary) []*PaypayPayinSummary {
	paypayPayinSummaryDTOs := make([]*PaypayPayinSummary, len(paypayPayinSummaries))
	for i, paypayPayinSummary := range paypayPayinSummaries {
		paypayPayinSummaryDTOs[i] = ToPaypayPayinSummaryDTO(paypayPayinSummary)
	}
	return paypayPayinSummaryDTOs
}

func ToPaypayPayinSummaryModels(paypayPayinSummaryDTOs []*PaypayPayinSummary) []*model.PaypayPayinSummary {
	paypayPayinSummaryModels := make([]*model.PaypayPayinSummary, len(paypayPayinSummaryDTOs))
	for i, paypayPayinSummaryDTO := range paypayPayinSummaryDTOs {
		paypayPayinSummaryModels[i] = paypayPayinSummaryDTO.ToPaypayPayinSummaryModel()
	}
	return paypayPayinSummaryModels
}
