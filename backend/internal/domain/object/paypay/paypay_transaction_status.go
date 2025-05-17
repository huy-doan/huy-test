package object

type PaypayTransactionStatus int

const (
	TransactionComplete        PaypayTransactionStatus = 1 // 取引完了
	TransactionAccepted        PaypayTransactionStatus = 2 // 取引受付完了
	RefundComplete             PaypayTransactionStatus = 3 // 返金完了
	TransactionCancelled       PaypayTransactionStatus = 4 // 取引取消
	TransactionAcceptCancelled PaypayTransactionStatus = 5 // 取引受付取消
	Adjustment                 PaypayTransactionStatus = 6 // 調整
	RemittanceComplete         PaypayTransactionStatus = 7 // 送金完了
)
