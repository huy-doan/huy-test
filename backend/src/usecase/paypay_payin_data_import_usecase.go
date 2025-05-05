package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
)

type DataImportUsecase struct {
	PayinDetailRepo      repositories.PayinDetailRepository
	PayinSummaryRepo     repositories.PayinSummaryRepository
	PayinTransactionRepo repositories.PayinTransactionRepository
	PayinFileUsecase     *PayinFileUsecase
}

func NewDataImportUsecase(
	payinDetailRepo repositories.PayinDetailRepository,
	payinSummaryRepo repositories.PayinSummaryRepository,
	payinTransactionRepo repositories.PayinTransactionRepository,
	payinFileUsecase *PayinFileUsecase,
) *DataImportUsecase {
	return &DataImportUsecase{
		PayinDetailRepo:      payinDetailRepo,
		PayinSummaryRepo:     payinSummaryRepo,
		PayinTransactionRepo: payinTransactionRepo,
		PayinFileUsecase:     payinFileUsecase,
	}
}

func (u *DataImportUsecase) InsertPayinDetailBatch(ctx context.Context, payinFileID int, records []map[string]string) error {
	var details []*models.PayPayPayinDetail
	for _, record := range records {
		transactionAmount, err := strconv.ParseFloat(record["transaction_amount"], 64)
		if err != nil {
			return err
		}
		refundAmount, _ := strconv.ParseFloat(record["refund_amount"], 64)
		usageFee, _ := strconv.ParseFloat(record["usage_fee"], 64)
		platformFee, _ := strconv.ParseFloat(record["platform_fee"], 64)
		initialFee, _ := strconv.ParseFloat(record["initial_fee"], 64)
		tax, _ := strconv.ParseFloat(record["tax"], 64)
		cashback, _ := strconv.ParseFloat(record["cashback"], 64)
		adjustment, _ := strconv.ParseFloat(record["adjustment"], 64)
		fee, _ := strconv.ParseFloat(record["fee"], 64)
		amount, _ := strconv.ParseFloat(record["amount"], 64)
		cutoffDate, _ := time.Parse(time.RFC3339, record["cutoff_date"])

		detail := &models.PayPayPayinDetail{
			PayinFileID:       payinFileID,
			PaymentMerchantID: record["payment_merchant_id"],
			StoreNumber:       record["store_number"],
			CutoffDate:        &cutoffDate,
			TransactionAmount: transactionAmount,
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
		details = append(details, detail)
	}
	return u.PayinDetailRepo.BulkInsert(ctx, details)
}

func (u *DataImportUsecase) InsertPayinSummaryBatch(ctx context.Context, payinFileID int, records []map[string]string) error {
	var summaries []*models.PayPayPayinSummary
	for _, record := range records {
		transactionAmount, err := strconv.ParseFloat(record["transaction_amount"], 64)
		if err != nil {
			return err
		}
		refundAmount, _ := strconv.ParseFloat(record["refund_amount"], 64)
		usageFee, _ := strconv.ParseFloat(record["usage_fee"], 64)
		platformFee, _ := strconv.ParseFloat(record["platform_fee"], 64)
		initialFee, _ := strconv.ParseFloat(record["initial_fee"], 64)
		tax, _ := strconv.ParseFloat(record["tax"], 64)
		cashback, _ := strconv.ParseFloat(record["cashback"], 64)
		adjustment, _ := strconv.ParseFloat(record["adjustment"], 64)
		fee, _ := strconv.ParseFloat(record["fee"], 64)
		amount, _ := strconv.ParseFloat(record["amount"], 64)
		cutoffDate, _ := time.Parse(time.RFC3339, record["cutoff_date"])
		paymentDate, _ := time.Parse(time.RFC3339, record["payment_date"])

		summary := &models.PayPayPayinSummary{
			PayinFileID:       payinFileID,
			CorporateName:     record["corporate_name"],
			CutoffDate:        &cutoffDate,
			PaymentDate:       &paymentDate,
			TransactionAmount: transactionAmount,
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
	return u.PayinSummaryRepo.BulkInsert(ctx, summaries)
}

func (u *DataImportUsecase) InsertPayinTransactionBatch(ctx context.Context, payinFile *models.PayinFile, records []map[string]string) error {
	var transactions []*models.PayPayPayinTransaction
	for _, record := range records {
		paymentTransactionID, err := strconv.Atoi(record["payment_transaction_id"])
		if err != nil {
			return err
		}
		transactionAmount, _ := strconv.ParseFloat(record["transaction_amount"], 64)
		paymentMethod, _ := strconv.Atoi(record["paypay_payment_method"])
		transactionStatus, _ := strconv.Atoi(record["payment_transaction_status"])
		transactionAt, _ := time.Parse(time.RFC3339, record["transaction_at"])

		transaction := &models.PayPayPayinTransaction{
			PayinFileID:              payinFile.ID,
			PaymentTransactionID:     &paymentTransactionID,
			PaymentMerchantID:        func(s string) *string { return &s }(record["payment_merchant_id"]),
			MerchantCode:             func(s string) *string { return &s }(record["merchant_code"]),
			ShopCode:                 func(s string) *string { return &s }(record["shop_code"]),
			ShopName:                 func(s string) *string { return &s }(record["shop_name"]),
			TerminalCode:             func(s string) *string { return &s }(record["terminal_code"]),
			PaymentTransactionStatus: &transactionStatus,
			TransactionAt:            &transactionAt,
			TransactionAmount:        &transactionAmount,
			ReceiptNumber:            func(s string) *string { return &s }(record["receipt_number"]),
			PayPayPaymentMethod:      &paymentMethod,
			SSID:                     func(s string) *string { return &s }(record["ssid"]),
			MerchantOrderID:          func(s string) *string { return &s }(record["merchant_order_id"]),
		}
		transactions = append(transactions, transaction)
	}
	return u.PayinTransactionRepo.BulkInsert(ctx, transactions)
}
