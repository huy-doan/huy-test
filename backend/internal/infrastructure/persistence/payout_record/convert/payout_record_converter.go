package convert

import (
	modelPayoutRecord "github.com/huydq/test/internal/domain/model/payout_record"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	objectPayout "github.com/huydq/test/internal/domain/object/payout"
	"github.com/huydq/test/internal/infrastructure/persistence/payout_record/dto"
)

// ToPayoutRecordDTO converts a PayoutRecord domain model to a PayoutRecord
func ToPayoutRecordDTO(record *modelPayoutRecord.PayoutRecord) *dto.PayoutRecord {
	if record == nil {
		return nil
	}

	return &dto.PayoutRecord{
		ID:                    record.ID,
		ShopID:                record.ShopID,
		PayoutID:              record.PayoutID,
		TransactionID:         record.TransactionID,
		BankName:              record.BankName,
		BankCode:              record.BankCode,
		BranchName:            record.BranchName,
		BranchCode:            record.BranchCode,
		BankAccountType:       int(record.BankAccountType),
		AccountNo:             record.AccountNo,
		AccountName:           record.AccountName,
		Amount:                record.Amount,
		TransferStatus:        int(record.TransferStatus),
		SendingDate:           record.SendingDate,
		AozoraTransferApplyNo: record.AozoraTransferApplyNo,
		TransferRequestedAt:   record.TransferRequestedAt,
		TransferExecutedAt:    record.TransferExecutedAt,
		TransferRequestError:  record.TransferRequestError,
		IdempotencyKey:        record.IdempotencyKey,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
			DeletedAt: record.DeletedAt,
		},
	}
}

// ToPayoutRecordModel converts a PayoutRecord to a PayoutRecord domain model
func ToPayoutRecordModel(dtoObj *dto.PayoutRecord) *modelPayoutRecord.PayoutRecord {
	if dtoObj == nil {
		return nil
	}

	return &modelPayoutRecord.PayoutRecord{
		ID:                    dtoObj.ID,
		ShopID:                dtoObj.ShopID,
		PayoutID:              dtoObj.PayoutID,
		TransactionID:         dtoObj.TransactionID,
		BankName:              dtoObj.BankName,
		BankCode:              dtoObj.BankCode,
		BranchName:            dtoObj.BranchName,
		BranchCode:            dtoObj.BranchCode,
		BankAccountType:       objectPayout.BankAccountType(dtoObj.BankAccountType),
		AccountNo:             dtoObj.AccountNo,
		AccountName:           dtoObj.AccountName,
		Amount:                dtoObj.Amount,
		TransferStatus:        objectPayout.TransferStatus(dtoObj.TransferStatus),
		SendingDate:           dtoObj.SendingDate,
		AozoraTransferApplyNo: dtoObj.AozoraTransferApplyNo,
		TransferRequestedAt:   dtoObj.TransferRequestedAt,
		TransferExecutedAt:    dtoObj.TransferExecutedAt,
		TransferRequestError:  dtoObj.TransferRequestError,
		IdempotencyKey:        dtoObj.IdempotencyKey,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
			DeletedAt: dtoObj.DeletedAt,
		},
	}
}

// ToPayoutRecordDTOs converts a list of PayoutRecord domain models to a list of PayoutRecordDTOs
func ToPayoutRecordDTOs(records []*modelPayoutRecord.PayoutRecord) []*dto.PayoutRecord {
	if records == nil {
		return nil
	}

	result := make([]*dto.PayoutRecord, len(records))
	for i, record := range records {
		result[i] = ToPayoutRecordDTO(record)
	}
	return result
}

// ToPayoutRecordModels converts a list of PayoutRecordDTOs to a list of PayoutRecord domain models
func ToPayoutRecordModels(dtos []*dto.PayoutRecord) []*modelPayoutRecord.PayoutRecord {
	if dtos == nil {
		return nil
	}

	result := make([]*modelPayoutRecord.PayoutRecord, len(dtos))
	for i, dto := range dtos {
		result[i] = ToPayoutRecordModel(dto)
	}
	return result
}
