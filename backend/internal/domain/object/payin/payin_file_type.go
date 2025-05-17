package object

type PayinFileType int

const (
	PayinFileTypePaymentSummary     PayinFileType = 0 // 入金レポート
	PayinFileTypePaymentDetail      PayinFileType = 1 // 入金明細
	PayinFileTypePaymentTransaction PayinFileType = 2 // 入金取引明細
)
