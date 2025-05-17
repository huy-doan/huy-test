package model

import (
	"time"

	util "github.com/huydq/test/internal/domain/object/basedatetime"
	object "github.com/huydq/test/internal/domain/object/payout"
)

// PayoutRecordParams contains parameters for creating a new PayoutRecord
type PayoutRecordParams struct {
	ID                    int
	ShopID                int
	PayoutID              int
	TransactionID         int
	BankName              string
	BankCode              string
	BranchName            string
	BranchCode            string
	BankAccountType       object.BankAccountType
	AccountNo             string
	AccountName           string
	Amount                float64
	TransferStatus        object.TransferStatus
	SendingDate           *time.Time
	AozoraTransferApplyNo string
	TransferRequestedAt   *time.Time
	TransferExecutedAt    *time.Time
	TransferRequestError  string
	IdempotencyKey        string
	util.BaseColumnTimestamp
}

// PayoutRecord represents a record of a payout transaction
type PayoutRecord struct {
	ID int
	util.BaseColumnTimestamp

	ShopID                int
	PayoutID              int
	TransactionID         int
	BankName              string
	BankCode              string
	BranchName            string
	BranchCode            string
	BankAccountType       object.BankAccountType
	AccountNo             string
	AccountName           string
	Amount                float64
	TransferStatus        object.TransferStatus
	SendingDate           *time.Time
	AozoraTransferApplyNo string
	TransferRequestedAt   *time.Time
	TransferExecutedAt    *time.Time
	TransferRequestError  string
	IdempotencyKey        string
}

// NewPayoutRecord creates a new payout record instance with the given parameters
func NewPayoutRecord(params PayoutRecordParams) *PayoutRecord {
	return &PayoutRecord{
		ID:                    params.ID,
		ShopID:                params.ShopID,
		PayoutID:              params.PayoutID,
		TransactionID:         params.TransactionID,
		BankName:              params.BankName,
		BankCode:              params.BankCode,
		BranchName:            params.BranchName,
		BranchCode:            params.BranchCode,
		BankAccountType:       params.BankAccountType,
		AccountNo:             params.AccountNo,
		AccountName:           params.AccountName,
		Amount:                params.Amount,
		TransferStatus:        params.TransferStatus,
		SendingDate:           params.SendingDate,
		AozoraTransferApplyNo: params.AozoraTransferApplyNo,
		TransferRequestedAt:   params.TransferRequestedAt,
		TransferExecutedAt:    params.TransferExecutedAt,
		TransferRequestError:  params.TransferRequestError,
		IdempotencyKey:        params.IdempotencyKey,
		BaseColumnTimestamp:   params.BaseColumnTimestamp,
	}
}

// MarkAsProcessed marks the transfer as processed
func (p *PayoutRecord) MarkAsProcessed(executedAt time.Time) {
	p.TransferStatus = object.TransferStatusProcessed
	p.TransferExecutedAt = &executedAt
}
