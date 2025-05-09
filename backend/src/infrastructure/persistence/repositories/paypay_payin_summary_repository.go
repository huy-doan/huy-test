package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"gorm.io/gorm"
)

type payinSummaryRepository struct {
	db *gorm.DB
}

func NewPayinSummaryRepository(db *gorm.DB) repositories.PayinSummaryRepository {
	return &payinSummaryRepository{db: db}
}

func (r *payinSummaryRepository) BulkInsert(ctx context.Context, summaries []*models.PayPayPayinSummary) error {
	return r.db.WithContext(ctx).Create(&summaries).Error
}
