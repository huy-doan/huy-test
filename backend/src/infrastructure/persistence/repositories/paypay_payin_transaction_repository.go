package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
	"gorm.io/gorm"
)

type payinTransactionRepository struct {
	db *gorm.DB
}

func NewPayinTransactionRepository(db *gorm.DB) repositories.PayinTransactionRepository {
	return &payinTransactionRepository{db: db}
}

func (r *payinTransactionRepository) BulkInsert(ctx context.Context, transactions []*models.PayPayPayinTransaction) error {
	return r.db.WithContext(ctx).Create(&transactions).Error
}
