package models

// PayinFileType constants for payin_file_type field
const (
	PayinFileTypePaymentReport = 1 // 入金レポート
	PayinFileTypePaymentDetail = 2 // 入金明細
)

// PayinFile represents the payin_file table
type PayinFile struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	PaymentProviderID    int    `json:"payment_provider_id"`
	PayinFileGroupID     *int   `json:"payin_file_group_id"`
	FileName             string `json:"file_name"`
	FileContentKey       string `json:"file_content_key"`
	PayinFileType        *int   `json:"payin_file_type"`
	HasDataRecord        bool   `json:"has_data_record"`
	AddedManually        bool   `json:"added_manually"`
	ContentAddedManually string `json:"content_added_manually"`
}

// TableName specifies the table name for PayinFile
func (PayinFile) TableName() string {
	return "payin_file"
}
