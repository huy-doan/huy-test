package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
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
