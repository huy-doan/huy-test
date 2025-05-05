package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
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
