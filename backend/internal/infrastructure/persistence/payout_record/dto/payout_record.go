package dto

import (
	"time"

	util "github.com/huydq/test/internal/domain/object/basedatetime"
	payoutDto "github.com/huydq/test/internal/infrastructure/persistence/payout/dto"
)

const (
	TransferStatusInProgress     int = 1 // 振込中
	TransferStatusWhitelistError int = 2 // ホワイトリスト追加エラー
	TransferStatusApiError       int = 3 // 振込依頼APIエラー
	TransferStatusFailed         int = 4 // 振込依頼失敗
	TransferStatusRequested      int = 5 // 振込依頼済み
	TransferStatusProcessed      int = 6 // 送金手続き済み

	BankAccountTypeOrdinary int = 1 // 普通預金
	BankAccountTypeCurrent  int = 2 // 当座預金
	BankAccountTypeFixed    int = 3 // 定期預金
)

type PayoutRecord struct {
	ID int `json:"id"`
	util.BaseColumnTimestamp

	ShopID                int        `json:"shop_id"`
	PayoutID              int        `json:"payout_id"`
	TransactionID         int        `json:"transaction_id"`
	BankName              string     `json:"bank_name"`
	BankCode              string     `json:"bank_code"`
	BranchName            string     `json:"branch_name"`
	BranchCode            string     `json:"branch_code"`
	BankAccountType       int        `json:"bank_account_type"`
	AccountNo             string     `json:"account_no"`
	AccountName           string     `json:"account_name"`
	Amount                float64    `json:"amount"`
	TransferStatus        int        `json:"transfer_status"`
	SendingDate           *time.Time `json:"sending_date"`
	AozoraTransferApplyNo string     `json:"aozora_transfer_apply_no"`
	TransferRequestedAt   *time.Time `json:"transfer_requested_at"`
	TransferExecutedAt    *time.Time `json:"transfer_executed_at"`
	TransferRequestError  string     `json:"transfer_request_error"`
	IdempotencyKey        string     `json:"idempotency_key"`

	// Shop        *Merchant    `json:"shop,omitempty" gorm:"foreignKey:ShopID"`
	Payout *payoutDto.Payout `json:"payout,omitempty" gorm:"foreignKey:PayoutID"`
	// Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
}

func (PayoutRecord) TableName() string {
	return "payout_record"
}
