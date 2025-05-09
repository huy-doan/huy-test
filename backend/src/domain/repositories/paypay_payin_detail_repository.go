package repositories

import (
	"context"

	"github.com/huydq/test/src/domain/models"
)

type PayinDetailRepository interface {
	BulkInsert(ctx context.Context, details []*models.PayPayPayinDetail) error
}
