package task

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

type InsertPayinSummaryTask struct {
	PaypayPayinSummaryRepo paypayRepo.PaypayPayinSummaryRepository
	appLogger              logger.Logger
}

func NewInsertPayinSummaryTask(repo paypayRepo.PaypayPayinSummaryRepository, appLogger logger.Logger) *InsertPayinSummaryTask {
	return &InsertPayinSummaryTask{
		PaypayPayinSummaryRepo: repo,
		appLogger:              appLogger,
	}
}

func parseFloatOrZero(s string) float64 {
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

func (t *InsertPayinSummaryTask) Do(ctx context.Context, payinFileID int, records []map[string]string) error {

	tx, err := database.GetTxOrDB(ctx)
	if err != nil {
		t.appLogger.ErrorWithContext("[InsertPayinSummaryTask] Error getting transaction: %v", err)
		return err
	}

	err = tx.Transaction(func(txCtx *gorm.DB) error {
		var summaries []*paypayModel.PaypayPayinSummary
		if len(records) > 0 {
			record := records[0]
			transactionAmount := parseFloatOrZero(record["transaction_amount"])
			refundAmount := parseFloatOrZero(record["refund_amount"])
			usageFee := parseFloatOrZero(record["usage_fee"])
			platformFee := parseFloatOrZero(record["platform_fee"])
			initialFee := parseFloatOrZero(record["initial_fee"])
			tax := parseFloatOrZero(record["tax"])
			cashback := parseFloatOrZero(record["cashback"])
			adjustment := parseFloatOrZero(record["adjustment"])
			fee := parseFloatOrZero(record["fee"])
			amount := parseFloatOrZero(record["amount"])

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

		if err := t.PaypayPayinSummaryRepo.BulkInsert(ctx, summaries); err != nil {
			log.Printf("[InsertPayinSummaryTask] Error inserting summaries: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		t.appLogger.ErrorWithContext("[InsertPayinSummaryTask] Error in transaction: %v", err)
		return err
	}
	return nil
}
