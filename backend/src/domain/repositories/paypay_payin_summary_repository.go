package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

type PayinSummaryRepository interface {
	BulkInsert(ctx context.Context, summaries []*models.PayPayPayinSummary) error
}
