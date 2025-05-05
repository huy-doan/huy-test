package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
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
