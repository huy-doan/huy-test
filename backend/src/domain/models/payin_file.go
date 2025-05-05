package models

// PayinFileType constants for payin_file_type field
const (
	PayinFileTypePaymentReport int = 1 // 入金レポート
	PayinFileTypePaymentDetail int = 2 // 入金明細
)

const (
	PaymentProviderID int = 1
)

const (
	StatusPending int = 1
	StatusSuccess int = 1
	StatusFailed  int = 2
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
	ImportStatus         int    `json:"import_status"`   // 0:pending, 1:success, 2:failed
	DownloadStatus       int    `json:"download_status"` // 0:pending, 1:success, 2:failed
	UploadStatus         int    `json:"upload_status"`   // 0:pending, 1:success, 2:failed
}

// TableName specifies the table name for PayinFile
func (PayinFile) TableName() string {
	return "payin_file"
}
