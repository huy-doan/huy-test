package paypay_batch

import (
	"context"

	repositories "github.com/huydq/test/internal/domain/repository/paypay_batch"
	payinFileModel "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_file/dto"
)

type PayinFileUsecase struct {
	repo repositories.PayinFileRepository
}

func NewPayinFileUsecase(repo repositories.PayinFileRepository) *PayinFileUsecase {
	return &PayinFileUsecase{repo: repo}
}

func (uc *PayinFileUsecase) CreateFile(ctx context.Context, file *payinFileModel.PayinFile) (*payinFileModel.PayinFile, error) {
	return uc.repo.Create(ctx, file)
}

func (uc *PayinFileUsecase) UpdateDownloadStatus(ctx context.Context, id int, status int) (*payinFileModel.PayinFile, error) {
	return uc.repo.UpdateStatus(ctx, id, "download_status", status)
}

func (uc *PayinFileUsecase) UpdateUploadStatus(ctx context.Context, id int, status int) (*payinFileModel.PayinFile, error) {
	return uc.repo.UpdateStatus(ctx, id, "upload_status", status)
}

func (uc *PayinFileUsecase) UpdateImportStatus(ctx context.Context, id int, status int) (*payinFileModel.PayinFile, error) {
	return uc.repo.UpdateStatus(ctx, id, "import_status", status)
}

func (uc *PayinFileUsecase) FileExistsAndDownloaded(ctx context.Context, filename string) (bool, error) {
	file, err := uc.repo.FindByFilename(ctx, filename)
	if err != nil || file.ID == 0 {
		return false, err
	}

	return file.DownloadStatus == payinFileModel.StatusSuccess, nil
}

func (uc *PayinFileUsecase) FindByFilename(ctx context.Context, filename string) (*payinFileModel.PayinFile, error) {
	payinFile, err := uc.repo.FindByFilename(ctx, filename)
	if err != nil {
		return nil, err
	}
	if payinFile == nil {
		return nil, nil // Return nil if no file is found
	}
	return payinFile, nil
}
