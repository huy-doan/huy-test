package usecase

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	paypayRepo "github.com/huydq/test/batch/domain/repository/paypay"
	paypayModel "github.com/huydq/test/internal/domain/model/paypay"
	"github.com/huydq/test/internal/pkg/database"
	"github.com/huydq/test/internal/pkg/logger"
	"gorm.io/gorm"
)

// PayinSummaryUsecase handles business logic for payin summaries
type PayinSummaryUsecase struct {
	repo      paypayRepo.PaypayPayinSummaryRepository
	appLogger logger.Logger
}

// NewPayinSummaryUsecase creates a new instance of PayinSummaryUsecase
func NewPayinSummaryUsecase(repo paypayRepo.PaypayPayinSummaryRepository, appLogger logger.Logger) *PayinSummaryUsecase {
	return &PayinSummaryUsecase{
		repo:      repo,
		appLogger: appLogger,
	}
}

// parseFloatOrZero parses a string to float64, returning 0 if parsing fails
func (uc *PayinSummaryUsecase) parseFloatOrZero(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// ProcessAndInsertSummaries processes summary records and inserts them into database
func (uc *PayinSummaryUsecase) ProcessAndInsertSummaries(ctx context.Context, payinFileID int, records []map[string]string) error {
	tx, err := database.GetTxOrDB(ctx)
	if err != nil {
		uc.appLogger.ErrorWithContext("[PayinSummaryUsecase] Error getting transaction: %v", err)
		return err
	}

	err = tx.Transaction(func(txCtx *gorm.DB) error {
		var summaries []*paypayModel.PaypayPayinSummary
		if len(records) > 0 {
			record := records[0]
			transactionAmount := uc.parseFloatOrZero(record["transaction_amount"])
			refundAmount := uc.parseFloatOrZero(record["refund_amount"])
			usageFee := uc.parseFloatOrZero(record["usage_fee"])
			platformFee := uc.parseFloatOrZero(record["platform_fee"])
			initialFee := uc.parseFloatOrZero(record["initial_fee"])
			tax := uc.parseFloatOrZero(record["tax"])
			cashback := uc.parseFloatOrZero(record["cashback"])
			adjustment := uc.parseFloatOrZero(record["adjustment"])
			fee := uc.parseFloatOrZero(record["fee"])
			amount := uc.parseFloatOrZero(record["amount"])

			// Defensive parse for cutoff_date
			var cutoffDatePtr *time.Time
			if s := strings.TrimSpace(record["cutoff_date"]); s != "" {
				if d, err := time.Parse("2006-01-02", s); err == nil {
					cutoffDatePtr = &d
				}
			}
			// Defensive parse for payment_date
			var paymentDatePtr *time.Time
			if s := strings.TrimSpace(record["payment_date"]); s != "" {
				if d, err := time.Parse("2006-01-02", s); err == nil {
					paymentDatePtr = &d
				}
			}

			summary := &paypayModel.PaypayPayinSummary{
				PayinFileID:       payinFileID,
				TransactionAmount: transactionAmount,
				CorporateName:     record["corporate_name"],
				CutoffDate:        cutoffDatePtr,
				PaymentDate:       paymentDatePtr,
				RefundAmount:      refundAmount,
				UsageFee:          usageFee,
				PlatformFee:       platformFee,
				InitialFee:        initialFee,
				Tax:               tax,
				Cashback:          cashback,
				Adjustment:        adjustment,
				Fee:               fee,
				Amount:            amount,
			}
			summaries = append(summaries, summary)
		}

		if err := uc.repo.BulkInsert(ctx, summaries); err != nil {
			log.Printf("[PayinSummaryUsecase] Error inserting summaries: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		uc.appLogger.ErrorWithContext("[PayinSummaryUsecase] Error in transaction: %v", err)
		return err
	}
	return nil
}