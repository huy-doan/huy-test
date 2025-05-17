package task

import (
	"context"
	"log"
	"time"

	payinUsecase "github.com/huydq/test/batch/usecase/payin"
	payinModel "github.com/huydq/test/internal/domain/model/payin"
	"github.com/huydq/test/internal/pkg/database"
	"gorm.io/gorm"
)

type CreateFileGroupTask struct {
	FileGroupUC *payinUsecase.PayinFileGroupUsecase
	ProviderID  int
}

func NewCreateFileGroupTask(
	fileGroupUC *payinUsecase.PayinFileGroupUsecase,
	providerID int,
) *CreateFileGroupTask {
	return &CreateFileGroupTask{
		FileGroupUC: fileGroupUC,
		ProviderID:  providerID,
	}
}

func (t *CreateFileGroupTask) Do(ctx context.Context) (int, error) {
	group := &payinModel.PayinFileGroup{
		FileGroupName:     time.Now().Format("20060102_150405"),
		PaymentProviderID: t.ProviderID,
		ImportTargetDate:  time.Now(),
	}

	tx, err := database.GetTxOrDB(ctx)
	if err != nil {
		return 0, err
	}

	err = tx.Transaction(func(txCtx *gorm.DB) error {
		err := t.FileGroupUC.CreateGroup(ctx, group)
		if err != nil {
			log.Printf("[CreateFileGroupTask] Failed to create group: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("[CreateFileGroupTask] Transaction error: %v", err)
		return 0, err
	}

	return group.ID, nil
}
