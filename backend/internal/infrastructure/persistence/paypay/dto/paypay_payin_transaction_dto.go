package dto

import (
	"time"

	paypayModel "github.com/huydq/test/internal/domain/model/paypay"
	paypayObject "github.com/huydq/test/internal/domain/object/paypay"
	"github.com/huydq/test/internal/infrastructure/persistence/payin/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

// PayPayPayinTransaction represents the paypay_payin_transaction table
type PaypayPayinTransaction struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`
	persistence.BaseColumnTimestamp

	PayinFileID          int                        `json:"payin_file_id"`
	PaymentTransactionID *string                    `json:"payment_transaction_id"`
	PaymentMerchantID    *string                    `json:"payment_merchant_id"`
	MerchantBusinessName *string                    `json:"merchant_business_name"`
	ShopID               *string                    `json:"shop_id"`
	ShopName             *string                    `json:"shop_name"`
	TerminalCode         *string                    `json:"terminal_code"`
	TransactionAt        *time.Time                 `json:"transaction_at"`
	TransactionAmount    *float64                   `json:"transaction_amount"`
	ReceiptNumber        *string                    `json:"receipt_number"`
	SSID                 *string                    `json:"ssid"`
	MerchantOrderID      *string                    `json:"merchant_order_id"`
	PaymentDetail        *paypayModel.PaymentDetail `json:"payment_detail"`

	PaymentTransactionStatus *paypayObject.PaypayTransactionStatus `json:"payment_transaction_status"`
	PaypayPaymentMethod      string                                `json:"paypay_payment_method"`

	PayinFile *dto.PayinFile `json:"payin_file,omitempty"`
}

// TableName specifies the table name for PayPayPayinTransaction
func (PaypayPayinTransaction) TableName() string {
	return "paypay_payin_transaction"
}

func (dto *PaypayPayinTransaction) ToPaypayPayinTransactionModel() *paypayModel.PaypayPayinTransaction {
	paypayPayinTransactionModel := &paypayModel.PaypayPayinTransaction{
		ID:                   dto.ID,
		PayinFileID:          dto.PayinFileID,
		PaymentTransactionID: dto.PaymentTransactionID,
		PaymentMerchantID:    dto.PaymentMerchantID,
		MerchantBusinessName: dto.MerchantBusinessName,
		ShopID:               dto.ShopID,
		ShopName:             dto.ShopName,
		TerminalCode:         dto.TerminalCode,
		TransactionAt:        dto.TransactionAt,
		TransactionAmount:    dto.TransactionAmount,
		ReceiptNumber:        dto.ReceiptNumber,
		SSID:                 dto.SSID,
		MerchantOrderID:      dto.MerchantOrderID,
		PaymentDetail:        dto.PaymentDetail,
	}
	paypayPayinTransactionModel.CreatedAt = dto.CreatedAt
	paypayPayinTransactionModel.UpdatedAt = dto.UpdatedAt
	return paypayPayinTransactionModel
}

func ToPaypayPayinTransactionDTO(paypayPayinTransactions *paypayModel.PaypayPayinTransaction) *PaypayPayinTransaction {
	paypayPayinTransactionDTO := &PaypayPayinTransaction{
		ID:                   paypayPayinTransactions.ID,
		PayinFileID:          paypayPayinTransactions.PayinFileID,
		PaymentTransactionID: paypayPayinTransactions.PaymentTransactionID,
		PaymentMerchantID:    paypayPayinTransactions.PaymentMerchantID,
		MerchantBusinessName: paypayPayinTransactions.MerchantBusinessName,
		ShopID:               paypayPayinTransactions.ShopID,
		ShopName:             paypayPayinTransactions.ShopName,
		TerminalCode:         paypayPayinTransactions.TerminalCode,
		TransactionAt:        paypayPayinTransactions.TransactionAt,
		TransactionAmount:    paypayPayinTransactions.TransactionAmount,
		ReceiptNumber:        paypayPayinTransactions.ReceiptNumber,
		SSID:                 paypayPayinTransactions.SSID,
		MerchantOrderID:      paypayPayinTransactions.MerchantOrderID,
		PaymentDetail:        paypayPayinTransactions.PaymentDetail,
	}
	paypayPayinTransactionDTO.CreatedAt = paypayPayinTransactions.CreatedAt
	paypayPayinTransactionDTO.UpdatedAt = paypayPayinTransactions.UpdatedAt
	return paypayPayinTransactionDTO
}

func ToPaypayPayinTransactionDTOs(paypayPayinTransactions []*paypayModel.PaypayPayinTransaction) []*PaypayPayinTransaction {
	paypayPayinTransactionDTOs := make([]*PaypayPayinTransaction, len(paypayPayinTransactions))
	for i, paypayPayinTransaction := range paypayPayinTransactions {
		paypayPayinTransactionDTOs[i] = ToPaypayPayinTransactionDTO(paypayPayinTransaction)
	}
	return paypayPayinTransactionDTOs
}

func ToPaypayPayinTransactionModels(paypayPayinTransactionDTOs []*PaypayPayinTransaction) []*paypayModel.PaypayPayinTransaction {
	paypayPayinTransactionModels := make([]*paypayModel.PaypayPayinTransaction, len(paypayPayinTransactionDTOs))
	for i, paypayPayinTransactionDTO := range paypayPayinTransactionDTOs {
		paypayPayinTransactionModels[i] = paypayPayinTransactionDTO.ToPaypayPayinTransactionModel()
	}
	return paypayPayinTransactionModels
}
