package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	paypayRepo "github.com/huydq/test/batch/domain/repository/paypay"
	paypayModel "github.com/huydq/test/internal/domain/model/paypay"
	paypayObject "github.com/huydq/test/internal/domain/object/paypay"
	"github.com/huydq/test/internal/pkg/database"
	"github.com/huydq/test/internal/pkg/logger"
	"gorm.io/gorm"
)

// PayinTransactionUsecase handles business logic for payin transactions
type PayinTransactionUsecase struct {
	repo      paypayRepo.PaypayPayinTransactionRepository
	appLogger logger.Logger
}

// NewPayinTransactionUsecase creates a new instance of PayinTransactionUsecase
func NewPayinTransactionUsecase(repo paypayRepo.PaypayPayinTransactionRepository, appLogger logger.Logger) *PayinTransactionUsecase {
	return &PayinTransactionUsecase{
		repo:      repo,
		appLogger: appLogger,
	}
}

// parseDatetime parses a datetime string into a *time.Time using multiple formats
func (uc *PayinTransactionUsecase) parseDatetime(val string) (*time.Time, error) {
	formats := []string{
		"2006/01/02 15:04:05",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return nil, nil
	}
	for _, f := range formats {
		t, err := time.Parse(f, val)
		if err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("cannot parse datetime: %s", val)
}

// ProcessAndInsertTransactions processes transaction records and inserts them into database
func (uc *PayinTransactionUsecase) ProcessAndInsertTransactions(ctx context.Context, payinFileID int, records []map[string]string) error {
	tx, err := database.GetTxOrDB(ctx)
	if err != nil {
		uc.appLogger.ErrorWithContext("[PayinTransactionUsecase] Error getting transaction: %v", err)
		return err
	}

	err = tx.Transaction(func(txCtx *gorm.DB) error {
		var transactions []*paypayModel.PaypayPayinTransaction
		for _, record := range records {

			// Helper to get pointer to string
			strPtr := func(s string) *string { s = strings.TrimSpace(s); return &s }

			// payment_transaction_id
			ptidStr := strings.TrimSpace(record["payment_transaction_id"])
			var paymentTransactionID *string
			if ptidStr != "" {
				paymentTransactionID = strPtr(ptidStr)
			}

			// transaction_amount
			transactionAmount, _ := strconv.ParseFloat(strings.TrimSpace(record["transaction_amount"]), 64)

			// transaction_status (string mapping)
			statusStr := strings.TrimSpace(record["payment_transaction_status"])
			var transactionStatus *paypayObject.PaypayTransactionStatus
			if statusStr != "" {
				switch statusStr {
				case "取引完了":
					val := paypayObject.TransactionComplete
					transactionStatus = &val
				case "取引受付完了":
					val := paypayObject.TransactionAccepted
					transactionStatus = &val
				case "返金完了":
					val := paypayObject.RefundComplete
					transactionStatus = &val
				case "取引取消":
					val := paypayObject.TransactionCancelled
					transactionStatus = &val
				case "取引受付取消":
					val := paypayObject.TransactionAcceptCancelled
					transactionStatus = &val
				case "調整":
					val := paypayObject.Adjustment
					transactionStatus = &val
				case "送金完了":
					val := paypayObject.RemittanceComplete
					transactionStatus = &val
				default:
					transactionStatus = nil
				}
			}
			transactionAt, _ := uc.parseDatetime(record["transaction_at"])
			transaction := &paypayModel.PaypayPayinTransaction{
				PayinFileID:              payinFileID,
				PaymentTransactionID:     paymentTransactionID,
				SSID:                     paymentTransactionID,
				PaymentMerchantID:        strPtr(record["payment_merchant_id"]),
				MerchantBusinessName:     strPtr(record["merchant_business_name"]),
				ShopID:                   strPtr(record["shop_id"]),
				ShopName:                 strPtr(record["shop_name"]),
				TerminalCode:             strPtr(record["terminal_code"]),
				PaymentTransactionStatus: transactionStatus,
				TransactionAt:            transactionAt,
				TransactionAmount:        &transactionAmount,
				ReceiptNumber:            strPtr(record["receipt_number"]),
				PaypayPaymentMethod:      record["paypay_payment_method"],
				MerchantOrderID:          strPtr(record["merchant_order_id"]),
			}
			transactions = append(transactions, transaction)
		}

		if len(transactions) == 0 {
			return nil
		}

		err := uc.repo.BulkInsert(ctx, transactions)
		if err != nil {
			uc.appLogger.ErrorWithContext("[PayinTransactionUsecase] Error inserting transactions: %v", err)
			return err
		}

		return nil
	})
	
	return err
}