package repositories

import (
	"context"

	repositories "github.com/huydq/test/internal/domain/repository/paypay_batch"
	summaryModels "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_summary/dto"
	"gorm.io/gorm"
)

type payinSummaryRepository struct {
	db *gorm.DB
}

func NewPayinSummaryRepository(db *gorm.DB) repositories.PayinSummaryRepository {
	return &payinSummaryRepository{db: db}
}

func (r *payinSummaryRepository) BulkInsert(ctx context.Context, summaries []*summaryModels.PayPayPayinSummary) error {
	return r.db.WithContext(ctx).Create(&summaries).Error
}
