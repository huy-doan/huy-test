package repositories

import (
	"context"

	models "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_transaction/dto"
)

type PayinTransactionRepository interface {
	BulkInsert(ctx context.Context, transactions []*models.PayPayPayinTransaction) error
}
