package usecase

import (
	"context"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
)

type PayinFileUsecase struct {
	repo repositories.PayinFileRepository
}

func NewPayinFileUsecase(repo repositories.PayinFileRepository) *PayinFileUsecase {
	return &PayinFileUsecase{repo: repo}
}

func (uc *PayinFileUsecase) CreateFile(ctx context.Context, file *models.PayinFile) error {
	return uc.repo.Create(ctx, file)
}

func (uc *PayinFileUsecase) UpdateDownloadStatus(ctx context.Context, id int, status int) error {
	return uc.repo.UpdateStatus(ctx, id, "download_status", status)
}

func (uc *PayinFileUsecase) UpdateUploadStatus(ctx context.Context, id int, status int) error {
	return uc.repo.UpdateStatus(ctx, id, "upload_status", status)
}

func (uc *PayinFileUsecase) UpdateImportStatus(ctx context.Context, id int, status int) error {
	return uc.repo.UpdateStatus(ctx, id, "import_status", status)
}

func (uc *PayinFileUsecase) FileExistsAndDownloaded(ctx context.Context, filename string) (bool, error) {
	id, err := uc.repo.FindIDByFilename(ctx, filename)
	if err != nil || id == 0 {
		return false, err
	}
	status, err := uc.repo.GetDownloadStatusByID(ctx, id)
	if err != nil {
		return false, err
	}
	return status == models.StatusSuccess, nil
}

func (uc *PayinFileUsecase) FindByFilename(ctx context.Context, filename string) (*models.PayinFile, error) {
	id, err := uc.repo.FindIDByFilename(ctx, filename)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, nil // Return nil if no file is found
	}

	payinFile, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return payinFile, nil
}
