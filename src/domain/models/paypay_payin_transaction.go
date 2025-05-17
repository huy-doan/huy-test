package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// PaymentTransactionStatus constants
const (
	TransactionComplete        int = 1 // 取引完了
	TransactionAccepted        int = 2 // 取引受付完了
	RefundComplete             int = 3 // 返金完了
	TransactionCancelled       int = 4 // 取引取消
	TransactionAcceptCancelled int = 5 // 取引受付取消
	Adjustment                 int = 6 // 調整
	RemittanceComplete         int = 7 // 送金完了
)

// PayPayPaymentMethod constants
const (
	PaymentMethodPayPayBalance     int = 1  // PayPay（残高）
	PaymentMethodCreditCard        int = 2  // クレジットカード
	PaymentMethodYahooMoney        int = 3  // Yahoo!マネー廃⽌
	PaymentMethodAlipay            int = 4  // Alipay
	PaymentMethodPayLater          int = 5  // あと払い（一括のみ）
	PaymentMethodPrepaidCode       int = 6  // プリペイドコード
	PaymentMethodLinePay           int = 7  // LinePay
	PaymentMethodPayPayCredit      int = 8  // PayPay（クレジット）
	PaymentMethodPayPayGiftCard    int = 9  // PayPay商品券
	PaymentMethodPayPayPoint       int = 10 // PayPayポイント
	PaymentMethodPayPayBankBalance int = 11 // PayPay銀行残高
)

// PaymentDetail represents payment method specific details
type PaymentDetail map[string]interface{}

// Value implements the driver.Valuer interface for database serialization
func (pd PaymentDetail) Value() (driver.Value, error) {
	return json.Marshal(pd)
}

// Scan implements the sql.Scanner interface for database deserialization
func (pd *PaymentDetail) Scan(value interface{}) error {
	if value == nil {
		*pd = make(PaymentDetail)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, pd)
}

// PayPayPayinTransaction represents the paypay_payin_transaction table
type PayPayPayinTransaction struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	PayinFileID              int            `json:"payin_file_id"`
	PaymentTransactionID     *int           `json:"payment_transaction_id"`
	PaymentMerchantID        *string        `json:"payment_merchant_id"`
	MerchantCode             *string        `json:"merchant_code"`
	ShopCode                 *string        `json:"shop_code"`
	ShopName                 *string        `json:"shop_name"`
	TerminalCode             *string        `json:"terminal_code"`
	PaymentTransactionStatus *int           `json:"payment_transaction_status"`
	TransactionAt            *time.Time     `json:"transaction_at"`
	TransactionAmount        *float64       `json:"transaction_amount"`
	ReceiptNumber            *string        `json:"receipt_number"`
	PayPayPaymentMethod      *int           `json:"paypay_payment_method"`
	SSID                     *string        `json:"ssid"`
	MerchantOrderID          *string        `json:"merchant_order_id"`
	PaymentDetail            *PaymentDetail `json:"payment_detail"`

	PayinFile *PayinFile `json:"payin_file,omitempty"`
}

// TableName specifies the table name for PayPayPayinTransaction
func (PayPayPayinTransaction) TableName() string {
	return "paypay_payin_transaction"
}
