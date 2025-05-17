package model

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// TransactionRecordType constants defined from the database comments
const (
	TransactionRecordTypeDeposit     = 1 // 入金
	TransactionRecordTypeFee         = 2 // 手数料
	TransactionRecordTypeTransferFee = 3 // 振込手数料
)

// TransactionRecord represents a transaction detail record entity in the system
type TransactionRecord struct {
	ID int `json:"id"`
	util.BaseColumnTimestamp

	TransactionID         int
	MerchantID            *int
	PayinDetailID         int
	PayinSummaryID        *int
	TransactionRecordType int
	Title                 string
	Amount                float64

	Transaction *Transaction
}
