package usecase

import (
	"context"

	"github.com/huydq/test/src/domain/models"
	"github.com/huydq/test/src/domain/repositories"
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
