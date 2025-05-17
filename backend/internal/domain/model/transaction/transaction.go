package model

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

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
	util.BaseColumnTimestamp

	ShopID            int
	TransactionStatus int
	PayoutID          int
	PayoutRecordID    int

	// Relationships - can be added later as needed
	TransactionRecords []TransactionRecord
}
