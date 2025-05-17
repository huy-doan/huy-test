package usecase

import (
	"context"

	repository "github.com/huydq/test/batch/domain/repository/payin"
	model "github.com/huydq/test/internal/domain/model/payin"
	object "github.com/huydq/test/internal/domain/object/payin"
)

type PayinFileUsecase struct {
	repo repository.PayinFileRepository
}

func NewPayinFileUsecase(repo repository.PayinFileRepository) *PayinFileUsecase {
	return &PayinFileUsecase{repo: repo}
}

func (uc *PayinFileUsecase) CreateFile(ctx context.Context, file *model.PayinFile) (*model.PayinFile, error) {
	return uc.repo.Create(ctx, file)
}

func (uc *PayinFileUsecase) UpdateDownloadStatus(ctx context.Context, file *model.PayinFile, status object.PayinFileStatus) error {
	if file.ID == 0 {
		return nil
	}
	file.UpdateDownloadStatus(status)
	return uc.repo.UpdateStatus(ctx, file)
}

func (uc *PayinFileUsecase) UpdateUploadStatus(ctx context.Context, file *model.PayinFile, status object.PayinFileStatus) error {
	if file.ID == 0 {
		return nil
	}
	file.UpdateUploadStatus(status)
	return uc.repo.UpdateStatus(ctx, file)
}

func (uc *PayinFileUsecase) UpdateImportStatus(ctx context.Context, file *model.PayinFile, status object.PayinFileStatus) error {
	if file.ID == 0 {
		return nil
	}
	file.UpdateImportStatus(status)
	return uc.repo.UpdateStatus(ctx, file)
}

func (uc *PayinFileUsecase) FileExistsAndDownloaded(ctx context.Context, filename string) (bool, error) {
	file, err := uc.repo.FindByFilename(ctx, filename)
	if err != nil || file == nil || file.ID == 0 {
		return false, err
	}
	return file.DownloadStatus == object.StatusSuccess, nil
}

func (uc *PayinFileUsecase) FindByFilename(ctx context.Context, filename string) (*model.PayinFile, error) {
	payinFile, err := uc.repo.FindByFilename(ctx, filename)
	if err != nil {
		return nil, err
	}
	if payinFile == nil {
		return nil, nil // Return nil if no file is found
	}
	return payinFile, nil
}
