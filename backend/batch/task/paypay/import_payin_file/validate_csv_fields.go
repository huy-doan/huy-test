package task

import object "github.com/huydq/test/internal/domain/object/payin"

// Required headers for each table
var RequiredPayinSummaryHeaders = []string{
	"payment_date",
	"transaction_amount",
	"fee",
	"amount",
	"cutoff_date",
	"refund_amount",
	"usage_fee",
	"platform_fee",
	"initial_fee",
	"tax",
	"cashback",
	"adjustment",
	"corporate_name",
}
var RequiredPayinDetailHeaders = []string{
	"payment_merchant_id",
	"merchant_business_name",
	"cutoff_date",
	"transaction_amount",
	"refund_amount",
	"usage_fee",
	"platform_fee",
	"initial_fee",
	"tax",
	"cashback",
	"adjustment",
	"fee",
	"amount",
}

var RequiredPayinTransactionHeaders = []string{"payment_transaction_id", "payment_merchant_id", "shop_id", "terminal_code", "payment_transaction_status", "receipt_number", "paypay_payment_method", "merchant_order_id", "merchant_business_name", "shop_name", "transaction_at", "transaction_amount", "payment_detail"}

// RequiredCSVHeaders maps file type to required headers (for backward compatibility)
var RequiredCSVHeaders = map[object.PayinFileType][]string{
	0: RequiredPayinSummaryHeaders,     // PayinFileTypePaymentSummary (first section)
	1: RequiredPayinDetailHeaders,      // PayinFileTypePaymentDetail (second section)
	2: RequiredPayinTransactionHeaders, // PayinFileTypePaymentTransaction (single-section)
}

type ValidateCSVFieldsTask struct{}

func NewValidateCSVFieldsTask() *ValidateCSVFieldsTask {
	return &ValidateCSVFieldsTask{}
}

// Do returns true if all required headers exist in the row, false if any is missing
func (t *ValidateCSVFieldsTask) Do(row map[string]string, fileType object.PayinFileType) bool {
	required, ok := RequiredCSVHeaders[fileType]
	if !ok {
		return false
	}
	for _, h := range required {
		if _, exists := row[h]; !exists {
			return false
		}
	}
	return true
}

// CheckHeaders checks if all required headers exist in the provided headers slice.
// If missing, you should update import status to failed in DB before any further processing.
func (t *ValidateCSVFieldsTask) CheckHeaders(headers []string, fileType object.PayinFileType) (bool, []string) {
	required, ok := RequiredCSVHeaders[fileType]
	if !ok {
		return false, required
	}
	headerSet := make(map[string]struct{}, len(headers))
	for _, h := range headers {
		headerSet[h] = struct{}{}
	}
	for _, req := range required {
		if _, exists := headerSet[req]; !exists {
			return false, required
		}
	}
	return true, required
}

// Validate all required headers for payin summary
func (t *ValidateCSVFieldsTask) CheckPayinSummaryHeaders(headers []string) (bool, []string) {
	headerSet := make(map[string]struct{}, len(headers))
	for _, h := range headers {
		headerSet[h] = struct{}{}
	}
	for _, req := range RequiredPayinSummaryHeaders {
		if _, exists := headerSet[req]; !exists {
			return false, RequiredPayinSummaryHeaders
		}
	}
	return true, RequiredPayinSummaryHeaders
}

// Validate all required headers for payin detail
func (t *ValidateCSVFieldsTask) CheckPayinDetailHeaders(headers []string) (bool, []string) {
	headerSet := make(map[string]struct{}, len(headers))
	for _, h := range headers {
		headerSet[h] = struct{}{}
	}
	for _, req := range RequiredPayinDetailHeaders {
		if _, exists := headerSet[req]; !exists {
			return false, RequiredPayinDetailHeaders
		}
	}
	return true, RequiredPayinDetailHeaders
}

// Validate all required headers for payin transaction
func (t *ValidateCSVFieldsTask) CheckPayinTransactionHeaders(headers []string) (bool, []string) {
	headerSet := make(map[string]struct{}, len(headers))
	for _, h := range headers {
		headerSet[h] = struct{}{}
	}
	for _, req := range RequiredPayinTransactionHeaders {
		if _, exists := headerSet[req]; !exists {
			return false, RequiredPayinTransactionHeaders
		}
	}
	return true, RequiredPayinTransactionHeaders
}
