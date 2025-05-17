package model

import (
	"time"

	util "github.com/huydq/test/internal/domain/object/basedatetime"
	object "github.com/huydq/test/internal/domain/object/payin"
)

const (
	PaymentProviderID int = 1
)

// PayinFile represents the payin_file table
type PayinFile struct {
	ID int
	util.BaseColumnTimestamp

	PaymentProviderID    int
	PayinFileGroupID     *int
	FileName             string
	FileContentKey       string
	PayinFileType        object.PayinFileType
	HasDataRecord        bool
	AddedManually        bool
	ContentAddedManually *string
	CreatedAt            time.Time
	ImportStatus         object.PayinFileStatus
	DownloadStatus       object.PayinFileStatus
	UploadStatus         object.PayinFileStatus
}

// UpdateDownloadStatus updates the download status
func (p *PayinFile) UpdateDownloadStatus(status object.PayinFileStatus) {
	p.DownloadStatus = status
}

// UpdateUploadStatus updates the upload status
func (p *PayinFile) UpdateUploadStatus(status object.PayinFileStatus) {
	p.UploadStatus = status
}

// UpdateImportStatus updates the import status
func (p *PayinFile) UpdateImportStatus(status object.PayinFileStatus) {
	p.ImportStatus = status
}
