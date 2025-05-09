package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
)

type PayinTransactionRepository interface {
	BulkInsert(ctx context.Context, transactions []*models.PayPayPayinTransaction) error
}
