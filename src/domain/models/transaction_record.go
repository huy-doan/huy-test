package models

// TransactionRecordType constants defined from the database comments
const (
	TransactionRecordTypeDeposit     = 1 // 入金
	TransactionRecordTypeFee         = 2 // 手数料
	TransactionRecordTypeTransferFee = 3 // 振込手数料
)

// TransactionRecord represents a transaction detail record entity in the system
type TransactionRecord struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	TransactionID         int     `json:"transaction_id"`
	MerchantID            *int    `json:"merchant_id"`
	PayinDetailID         int     `json:"payin_detail_id"`
	PayinSummaryID        *int    `json:"payin_summary_id"`
	TransactionRecordType int     `json:"transaction_record_type"`
	Title                 string  `json:"title"`
	Amount                float64 `json:"amount"`

	// Relationships - can be added later as needed
	Transaction *Transaction `json:"transaction,omitempty"`
}

// TableName specifies the database table name
func (TransactionRecord) TableName() string {
	return "transaction_record"
}
