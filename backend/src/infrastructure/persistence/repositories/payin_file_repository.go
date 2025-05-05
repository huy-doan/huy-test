package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/vnlab/makeshop-payment/src/domain/models"
	"github.com/vnlab/makeshop-payment/src/domain/repositories"
)

type payinFileRepo struct {
	db *gorm.DB
}

func NewPayinFileRepository(db *gorm.DB) repositories.PayinFileRepository {
	return &payinFileRepo{db: db}
}

func (r *payinFileRepo) Create(ctx context.Context, file *models.PayinFile) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *payinFileRepo) UpdateStatus(ctx context.Context, id int, field string, status int) error {
	return r.db.WithContext(ctx).
		Model(&models.PayinFile{}).
		Where("id = ?", id).
		Update(field, status).Error
}

func (r *payinFileRepo) FindIDByFilename(ctx context.Context, filename string) (int, error) {
	var payinFile models.PayinFile
	err := r.db.WithContext(ctx).
		Where("file_name = ?", filename).
		First(&payinFile).Error
	if err != nil {
		return 0, err
	}
	return payinFile.ID, nil
}

func (r *payinFileRepo) GetDownloadStatusByID(ctx context.Context, id int) (int, error) {
	var status int
	err := r.db.WithContext(ctx).
		Model(&models.PayinFile{}).
		Where("id = ?", id).
		Select("download_status").
		Scan(&status).Error
	return status, err
}

func (r *payinFileRepo) GetByID(ctx context.Context, id int) (*models.PayinFile, error) {
	var payinFile models.PayinFile
	if err := r.db.WithContext(ctx).First(&payinFile, id).Error; err != nil {
		return nil, err
	}
	return &payinFile, nil
}
