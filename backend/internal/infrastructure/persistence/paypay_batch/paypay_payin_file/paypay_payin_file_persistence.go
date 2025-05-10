package repositories

import (
	"context"

	"gorm.io/gorm"

	repositories "github.com/huydq/test/internal/domain/repository/paypay_batch"
	models "github.com/huydq/test/internal/infrastructure/persistence/paypay_batch/paypay_payin_file/dto"
)

type payinFileRepo struct {
	db *gorm.DB
}

func NewPayinFileRepository(db *gorm.DB) repositories.PayinFileRepository {
	return &payinFileRepo{db: db}
}

func (r *payinFileRepo) Create(ctx context.Context, file *models.PayinFile) (*models.PayinFile, error) {
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (r *payinFileRepo) UpdateStatus(ctx context.Context, id int, field string, status int) (*models.PayinFile, error) {
	if err := r.db.WithContext(ctx).
		Model(&models.PayinFile{}).
		Where("id = ?", id).
		Update(field, status).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *payinFileRepo) FindByFilename(ctx context.Context, filename string) (*models.PayinFile, error) {
	var file models.PayinFile
	if err := r.db.WithContext(ctx).
		Where("file_name = ?", filename).
		First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *payinFileRepo) GetByID(ctx context.Context, id int) (*models.PayinFile, error) {
	var file models.PayinFile
	if err := r.db.WithContext(ctx).First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}
