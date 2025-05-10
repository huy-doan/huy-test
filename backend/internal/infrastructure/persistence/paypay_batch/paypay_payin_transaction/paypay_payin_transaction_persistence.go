package repositories

import (
	"context"

	repositories "github.com/huydq/test/internal/domain/repository/paypay_batch"
	models "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_transaction/dto"
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
