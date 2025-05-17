package usecase

import (
	"context"
	"strconv"
	"strings"
	"time"

	paypayRepo "github.com/huydq/test/batch/domain/repository/paypay"
	paypayModel "github.com/huydq/test/internal/domain/model/paypay"
	"github.com/huydq/test/internal/pkg/database"
	"github.com/huydq/test/internal/pkg/logger"
	"gorm.io/gorm"
)

// PayinDetailUsecase handles business logic for payin details
type PayinDetailUsecase struct {
	repo      paypayRepo.PaypayPayinDetailRepository
	appLogger logger.Logger
}

func NewPayinDetailUsecase(repo paypayRepo.PaypayPayinDetailRepository, appLogger logger.Logger) *PayinDetailUsecase {
	return &PayinDetailUsecase{
		repo:      repo,
		appLogger: appLogger,
	}
}

// ProcessAndInsertDetails processes detail records and inserts them into database
func (uc *PayinDetailUsecase) ProcessAndInsertDetails(ctx context.Context, payinFileID int, records []map[string]string) error {
	tx, err := database.GetTxOrDB(ctx)
	if err != nil {
		uc.appLogger.ErrorWithContext("[PayinDetailUsecase] Error getting transaction: %v", err)
		return err
	}

	err = tx.Transaction(func(txCtx *gorm.DB) error {
		var details []*paypayModel.PaypayPayinDetail
		for _, record := range records {
			// Parse transaction amount
			transactionAmount, err := strconv.ParseFloat(record["transaction_amount"], 64)
			if err != nil {
				uc.appLogger.ErrorWithContext("[PayinDetailUsecase] Error parsing transaction_amount: %v", err)
				return err
			}

			// Parse other numeric fields
			refundAmount, _ := strconv.ParseFloat(record["refund_amount"], 64)
			usageFee, _ := strconv.ParseFloat(record["usage_fee"], 64)
			platformFee, _ := strconv.ParseFloat(record["platform_fee"], 64)
			initialFee, _ := strconv.ParseFloat(record["initial_fee"], 64)
			tax, _ := strconv.ParseFloat(record["tax"], 64)
			cashback, _ := strconv.ParseFloat(record["cashback"], 64)
			adjustment, _ := strconv.ParseFloat(record["adjustment"], 64)
			fee, _ := strconv.ParseFloat(record["fee"], 64)
			amount, _ := strconv.ParseFloat(record["amount"], 64)

			// Defensive parse for cutoff_date
			var cutoffDatePtr *time.Time
			if s := strings.TrimSpace(record["cutoff_date"]); s != "" {
				if d, err := time.Parse("2006-01-02", s); err == nil {
					cutoffDatePtr = &d
				}
			}

			// Create detail record
			detail := &paypayModel.PaypayPayinDetail{
				PayinFileID:          payinFileID,
				PaymentMerchantID:    record["payment_merchant_id"],
				MerchantBusinessName: record["merchant_business_name"],
				CutoffDate:           cutoffDatePtr,
				TransactionAmount:    transactionAmount,
				RefundAmount:         refundAmount,
				UsageFee:             usageFee,
				PlatformFee:          platformFee,
				InitialFee:           initialFee,
				Tax:                  tax,
				Cashback:             cashback,
				Adjustment:           adjustment,
				Fee:                  fee,
				Amount:               amount,
			}
			details = append(details, detail)
		}

		// Insert all records in bulk
		err := uc.repo.BulkInsert(ctx, details)
		if err != nil {
			uc.appLogger.ErrorWithContext("[PayinDetailUsecase] Error inserting details: %v", err)
			return err
		}
		return nil
	})
	
	return err
}
