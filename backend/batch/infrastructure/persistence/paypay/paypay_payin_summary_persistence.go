package persistence

import (
	"context"

	repository "github.com/huydq/test/batch/domain/repository/paypay"
	model "github.com/huydq/test/internal/domain/model/paypay"
	dto "github.com/huydq/test/internal/infrastructure/persistence/paypay/dto"

	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

type PaypayPayinSummaryPersistence struct {
	db *gorm.DB
}

func NewPayinSummaryRepository(db *gorm.DB) repository.PaypayPayinSummaryRepository {
	return &PaypayPayinSummaryPersistence{db: db}
}

func (r *PaypayPayinSummaryPersistence) BulkInsert(ctx context.Context, summaries []*model.PaypayPayinSummary) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	paypayPayinSummaryDTOs := dto.ToPaypayPayinSummaryDTOs(summaries)
	return db.WithContext(ctx).Create(&paypayPayinSummaryDTOs).Error
}
