package repositories

import (
	"context"

	repositories "github.com/huydq/test/internal/domain/repository/paypay_batch"
	models "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_detail/dto"
	"gorm.io/gorm"
)

type payinDetailRepository struct {
	db *gorm.DB
}

func NewPayinDetailRepository(db *gorm.DB) repositories.PayinDetailRepository {
	return &payinDetailRepository{db: db}
}

func (r *payinDetailRepository) BulkInsert(ctx context.Context, details []*models.PayPayPayinDetail) error {
	return r.db.WithContext(ctx).Create(&details).Error
}
