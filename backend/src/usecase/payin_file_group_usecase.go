package usecase

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
)

type PayinFileGroupUsecase struct {
	repo repositories.PayinFileGroupRepository
}

func NewPayinFileGroupUsecase(repo repositories.PayinFileGroupRepository) *PayinFileGroupUsecase {
	return &PayinFileGroupUsecase{repo: repo}
}

func (uc *PayinFileGroupUsecase) CreateGroup(ctx context.Context, group *models.PayinFileGroup) error {
	return uc.repo.Create(ctx, group)
}
