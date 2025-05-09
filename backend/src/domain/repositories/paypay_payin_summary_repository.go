package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
)

type PayinSummaryRepository interface {
	BulkInsert(ctx context.Context, summaries []*models.PayPayPayinSummary) error
}
