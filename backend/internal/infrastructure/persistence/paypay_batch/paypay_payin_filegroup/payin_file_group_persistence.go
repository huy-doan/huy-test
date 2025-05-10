package repositories

import (
	"context"

	"gorm.io/gorm"

	repositories "github.com/huydq/test/internal/domain/repository/paypay_batch"
	models "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_filegroup/dto"
)

type payinFileGroupRepo struct {
	db *gorm.DB
}

func NewPayinFileGroupRepository(db *gorm.DB) repositories.PayinFileGroupRepository {
	return &payinFileGroupRepo{db: db}
}

func (r *payinFileGroupRepo) Create(ctx context.Context, group *models.PayinFileGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}
