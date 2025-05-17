package persistence

import (
	"context"

	repository "github.com/huydq/test/batch/domain/repository/paypay"
	model "github.com/huydq/test/internal/domain/model/paypay"
	dto "github.com/huydq/test/internal/infrastructure/persistence/paypay/dto"

	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

type PaypayPayinTransactionPersistence struct {
	db *gorm.DB
}

func NewPayinTransactionRepository(db *gorm.DB) repository.PaypayPayinTransactionRepository {
	return &PaypayPayinTransactionPersistence{db: db}
}

func (r *PaypayPayinTransactionPersistence) BulkInsert(ctx context.Context, transactions []*model.PaypayPayinTransaction) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	paypayPayinTransactionDTOs := dto.ToPaypayPayinTransactionDTOs(transactions)
	return db.WithContext(ctx).Create(&paypayPayinTransactionDTOs).Error
}
