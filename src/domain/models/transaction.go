package models

// Transaction status constants defined from the database comments
const (
	TransactionStatusProcessing        = 1 // 処理中
	TransactionStatusPendingApproval   = 2 // 承認待ち
	TransactionStatusApproved          = 3 // 承認済み
	TransactionStatusTransferRequested = 4 // 送金依頼済
	TransactionStatusTransferFailed    = 5 // 送金依頼失敗
)

// Transaction represents a transaction entity in the system
type Transaction struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	ShopID            int `json:"shop_id"`
	TransactionStatus int `json:"transaction_status"`
	PayoutID          int `json:"payout_id"`
	PayoutRecordID    int `json:"payout_record_id"`

	// Relationships - can be added later as needed
	TransactionRecords []TransactionRecord `json:"transaction_records,omitempty"`
}

// TableName specifies the database table name
func (Transaction) TableName() string {
	return "transaction"
}
