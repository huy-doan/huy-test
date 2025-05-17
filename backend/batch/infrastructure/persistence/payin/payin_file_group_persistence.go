package persistence

import (
	"context"

	"gorm.io/gorm"

	repository "github.com/huydq/test/batch/domain/repository/payin"
	model "github.com/huydq/test/internal/domain/model/payin"
	"github.com/huydq/test/internal/infrastructure/persistence/payin/dto"
	"github.com/huydq/test/internal/pkg/database"
)

type PayinFileGroupPersistence struct {
	db *gorm.DB
}

func NewPayinFileGroupRepository(db *gorm.DB) repository.PayinFileGroupRepository {
	return &PayinFileGroupPersistence{db: db}
}

func (r *PayinFileGroupPersistence) Create(ctx context.Context, fileGroup *model.PayinFileGroup) error {
	db, err := database.GetTxOrDB(ctx)
	if err != nil {
		return err
	}

	payinFileGroupDTO := dto.ToPayinFileGroupDTO(fileGroup)
	return db.Create(payinFileGroupDTO).Error
}
