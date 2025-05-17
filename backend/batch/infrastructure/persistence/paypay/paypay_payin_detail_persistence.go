package persistence

import (
	"context"

	repository "github.com/huydq/test/batch/domain/repository/paypay"
	model "github.com/huydq/test/internal/domain/model/paypay"
	dto "github.com/huydq/test/internal/infrastructure/persistence/paypay/dto"

	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

type PaypayPayinDetailPersistence struct {
	db *gorm.DB
}

func NewPayinDetailRepository(db *gorm.DB) repository.PaypayPayinDetailRepository {
	return &PaypayPayinDetailPersistence{db: db}
}

func (r *PaypayPayinDetailPersistence) BulkInsert(ctx context.Context, details []*model.PaypayPayinDetail) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	paypayPayinDetailDTOs := dto.ToPaypayPayinDetailDTOs(details)
	return db.WithContext(ctx).Create(&paypayPayinDetailDTOs).Error
}
