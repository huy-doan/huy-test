package dto

import (
	"time"

	payinFileGroup "github.com/huydq/test/internal/domain/model/payin"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
)

// PayinFileGroup represents the payin_file_group table
type PayinFileGroup struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`
	persistence.BaseColumnTimestamp

	FileGroupName     string     `json:"file_group_name"`
	PaymentProviderID int        `json:"payment_provider_id"`
	ImportTargetDate  time.Time  `json:"import_target_date"`
	ImportedAt        *time.Time `json:"imported_at"`
}

// TableName specifies the table name for PayinFileGroup
func (PayinFileGroup) TableName() string {
	return "payin_file_group"
}

func (dto *PayinFileGroup) ToPayinFileGroupModel() *payinFileGroup.PayinFileGroup {
	payinFileGroupModel := &payinFileGroup.PayinFileGroup{
		ID:                dto.ID,
		PaymentProviderID: dto.PaymentProviderID,
		FileGroupName:     dto.FileGroupName,
		ImportTargetDate:  dto.ImportTargetDate,
	}
	payinFileGroupModel.CreatedAt = dto.CreatedAt
	payinFileGroupModel.UpdatedAt = dto.UpdatedAt
	return payinFileGroupModel
}

func ToPayinFileGroupDTO(pfg *payinFileGroup.PayinFileGroup) *PayinFileGroup {
	payinFileGroupDTO := &PayinFileGroup{
		ID:                pfg.ID,
		PaymentProviderID: pfg.PaymentProviderID,
		FileGroupName:     pfg.FileGroupName,
		ImportTargetDate:  pfg.ImportTargetDate,
		ImportedAt:        pfg.ImportedAt,
	}
	payinFileGroupDTO.CreatedAt = pfg.CreatedAt
	payinFileGroupDTO.UpdatedAt = pfg.UpdatedAt
	return payinFileGroupDTO
}
