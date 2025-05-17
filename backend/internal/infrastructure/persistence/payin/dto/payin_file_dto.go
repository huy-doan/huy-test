package dto

import (
	payinFile "github.com/huydq/test/internal/domain/model/payin"
	object "github.com/huydq/test/internal/domain/object/payin"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

// PayinFile represents the payin_file table
type PayinFile struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`
	persistence.BaseColumnTimestamp

	PaymentProviderID    int     `json:"payment_provider_id"`
	PayinFileGroupID     *int    `json:"payin_file_group_id"`
	FileName             string  `json:"file_name"`
	FileContentKey       string  `json:"file_content_key"`
	HasDataRecord        bool    `json:"has_data_record"`
	AddedManually        bool    `json:"added_manually"`
	ContentAddedManually *string `json:"content_added_manually"`

	PayinFileType object.PayinFileType `json:"payin_file_type"`

	ImportStatus   object.PayinFileStatus `json:"import_status"`   // 0:pending, 1:success, 2:failed
	DownloadStatus object.PayinFileStatus `json:"download_status"` // 0:pending, 1:success, 2:failed
	UploadStatus   object.PayinFileStatus `json:"upload_status"`   // 0:pending, 1:success, 2:failed
}

// TableName specifies the table name for PayinFile
func (PayinFile) TableName() string {
	return "payin_file"
}

func (dto *PayinFile) ToPayinFileModel() *payinFile.PayinFile {
	payinFileModel := &payinFile.PayinFile{
		ID:                   dto.ID,
		PaymentProviderID:    dto.PaymentProviderID,
		PayinFileGroupID:     dto.PayinFileGroupID,
		FileName:             dto.FileName,
		FileContentKey:       dto.FileContentKey,
		HasDataRecord:        dto.HasDataRecord,
		AddedManually:        dto.AddedManually,
		ContentAddedManually: dto.ContentAddedManually,
		PayinFileType:        dto.PayinFileType,
		ImportStatus:         dto.ImportStatus,
		DownloadStatus:       dto.DownloadStatus,
		UploadStatus:         dto.UploadStatus,
	}
	payinFileModel.CreatedAt = dto.CreatedAt
	payinFileModel.UpdatedAt = dto.UpdatedAt
	return payinFileModel
}

func ToPayinFileDTO(pf *payinFile.PayinFile) *PayinFile {
	payinFileDTO := &PayinFile{
		ID:                   pf.ID,
		PaymentProviderID:    pf.PaymentProviderID,
		PayinFileGroupID:     pf.PayinFileGroupID,
		FileName:             pf.FileName,
		FileContentKey:       pf.FileContentKey,
		HasDataRecord:        pf.HasDataRecord,
		AddedManually:        pf.AddedManually,
		ContentAddedManually: pf.ContentAddedManually,
		PayinFileType:        pf.PayinFileType,
		ImportStatus:         pf.ImportStatus,
		DownloadStatus:       pf.DownloadStatus,
		UploadStatus:         pf.UploadStatus,
	}
	payinFileDTO.CreatedAt = pf.CreatedAt
	payinFileDTO.UpdatedAt = pf.UpdatedAt
	return payinFileDTO
}
