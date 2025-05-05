package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

type PayinTransactionRepository interface {
	BulkInsert(ctx context.Context, transactions []*models.PayPayPayinTransaction) error
}
