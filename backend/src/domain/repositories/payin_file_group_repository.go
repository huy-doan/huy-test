package repositories

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
)

type PayinFileGroupRepository interface {
	//Create methods for PayinFileGroupRepository
	Create(ctx context.Context, group *models.PayinFileGroup) error
}
