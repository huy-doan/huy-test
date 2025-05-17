package object

type PaypayPaymentMethod int

// PayPayPaymentMethod constants
const (
	PaymentMethodPayPayBalance     PaypayPaymentMethod = 1  // PayPay（残高）
	PaymentMethodCreditCard        PaypayPaymentMethod = 2  // クレジットカード
	PaymentMethodYahooMoney        PaypayPaymentMethod = 3  // Yahoo!マネー廃⽌
	PaymentMethodAlipay            PaypayPaymentMethod = 4  // Alipay
	PaymentMethodPayLater          PaypayPaymentMethod = 5  // あと払い（一括のみ）
	PaymentMethodPrepaidCode       PaypayPaymentMethod = 6  // プリペイドコード
	PaymentMethodLinePay           PaypayPaymentMethod = 7  // LinePay
	PaymentMethodPayPayCredit      PaypayPaymentMethod = 8  // PayPay（クレジット）
	PaymentMethodPayPayGiftCard    PaypayPaymentMethod = 9  // PayPay商品券
	PaymentMethodPayPayPoint       PaypayPaymentMethod = 10 // PayPayポイント
	PaymentMethodPayPayBankBalance PaypayPaymentMethod = 11 // PayPay銀行残高
)
