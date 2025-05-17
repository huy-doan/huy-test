package repository

import (
	"context"

	model "github.com/huydq/test/internal/domain/model/payin"
)

type PayinFileGroupRepository interface {
	//Create methods for PayinFileGroupRepository
	Create(ctx context.Context, group *model.PayinFileGroup) error
}
