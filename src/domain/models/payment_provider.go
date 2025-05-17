package models

// PaymentProvider represents the payment_provider table
type PaymentProvider struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	Code string `json:"code"`
	Name string `json:"name"`
}

// TableName specifies the table name for PaymentProvider
func (PaymentProvider) TableName() string {
	return "payment_provider"
}
