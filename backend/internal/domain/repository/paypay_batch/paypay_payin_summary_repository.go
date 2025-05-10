package repositories

import (
	"context"

	models "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_summary/dto"
)

type PayinSummaryRepository interface {
	BulkInsert(ctx context.Context, summaries []*models.PayPayPayinSummary) error
}
