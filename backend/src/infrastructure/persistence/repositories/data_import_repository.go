package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"gorm.io/gorm"
)

type DataImportRepository struct {
	db *gorm.DB
}

func NewDataImportRepository(db *gorm.DB) *DataImportRepository {
	return &DataImportRepository{db: db}
}

func (r *DataImportRepository) BulkInsert(ctx context.Context, transactions []*models.PayPayPayinTransaction) error {
	return r.db.WithContext(ctx).Create(&transactions).Error
}
