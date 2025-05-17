package persistence

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	repository "github.com/huydq/test/batch/domain/repository/payin"
	model "github.com/huydq/test/internal/domain/model/payin"
	"github.com/huydq/test/internal/infrastructure/persistence/payin/dto"
	"github.com/huydq/test/internal/pkg/database"
)

type PayinFilePersistence struct {
	db *gorm.DB
}

func NewPayinFileRepository(db *gorm.DB) repository.PayinFileRepository {
	return &PayinFilePersistence{db: db}
}

func (r *PayinFilePersistence) Create(ctx context.Context, file *model.PayinFile) (*model.PayinFile, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	payinFileDTO := dto.ToPayinFileDTO(file)
	if err := db.Create(&payinFileDTO).Error; err != nil {
		return nil, err
	}

	fileModel := payinFileDTO.ToPayinFileModel()
	return fileModel, nil
}

func (r *PayinFilePersistence) UpdateStatus(ctx context.Context, file *model.PayinFile) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	return db.Model(&model.PayinFile{}).
		Select("download_status", "upload_status", "import_status").
		Where("id = ?", file.ID).
		Updates(file).Error
}

func (r *PayinFilePersistence) FindByFilename(ctx context.Context, filename string) (*model.PayinFile, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		log.Printf("[FindByFilename] Failed to get DB: %v", err)
		return nil, err
	}

	var fileDTO dto.PayinFile
	if err := db.Where("file_name = ?", filename).First(&fileDTO).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return fileDTO.ToPayinFileModel(), nil
}

func (r *PayinFilePersistence) GetByID(ctx context.Context, id int) (*model.PayinFile, error) {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return nil, err
	}

	var fileDTO dto.PayinFile
	if err := db.First(&fileDTO, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return fileDTO.ToPayinFileModel(), nil
}
