package model

import (
	"time"

	model "github.com/huydq/test/internal/domain/model/payin"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	paypayObject "github.com/huydq/test/internal/domain/object/paypay"
)

// PayPayPayinTransaction represents the paypay_payin_transaction table
type PaypayPayinTransaction struct {
	ID int
	util.BaseColumnTimestamp

	PayinFileID              int
	PaymentTransactionID     *string
	PaymentMerchantID        *string
	MerchantBusinessName     *string
	ShopID                   *string
	ShopName                 *string
	TerminalCode             *string
	PaymentTransactionStatus *paypayObject.PaypayTransactionStatus
	TransactionAt            *time.Time
	TransactionAmount        *float64
	ReceiptNumber            *string
	PaypayPaymentMethod      string
	SSID                     *string
	MerchantOrderID          *string

	// TODO: PaymentDetail 型の正確なJSON構造がまだ判明していないため、定義ができていません。
	// 詳細な構造が把握できれば、map[string]any の代わりに適切な構造体を作成する方が良いでしょう。
	PaymentDetail *PaymentDetail

	PayinFile *model.PayinFile
}
