package convert

import (
	modelPayout "github.com/huydq/test/internal/domain/model/payout"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	objectPayout "github.com/huydq/test/internal/domain/object/payout"
	"github.com/huydq/test/internal/infrastructure/persistence/payout/dto"
	userDto "github.com/huydq/test/internal/infrastructure/persistence/user/dto"
)

// ToPayoutDTO converts a Payout domain model to a PayoutDTO
func ToPayoutDTO(payout *modelPayout.Payout) *dto.PayoutDTO {
	if payout == nil {
		return nil
	}

	return &dto.PayoutDTO{
		ID:                    payout.ID,
		PayoutStatus:          int(payout.PayoutStatus),
		Total:                 payout.Total,
		TotalCount:            payout.TotalCount,
		SendingDate:           payout.SendingDate,
		SentDate:              payout.SentDate,
		AozoraTransferApplyNo: payout.AozoraTransferApplyNo,
		ApprovalID:            payout.ApprovalID,
		UserID:                payout.UserID,
		User:                  userDto.ToUserDTO(payout.User),
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: payout.CreatedAt,
			UpdatedAt: payout.UpdatedAt,
			DeletedAt: payout.DeletedAt,
		},
	}
}

func ToPayoutModel(dtoObj *dto.PayoutDTO) *modelPayout.Payout {
	if dtoObj == nil {
		return nil
	}

	return &modelPayout.Payout{
		ID:                    dtoObj.ID,
		PayoutStatus:          objectPayout.PayoutStatus(dtoObj.PayoutStatus),
		Total:                 dtoObj.Total,
		TotalCount:            dtoObj.TotalCount,
		SendingDate:           dtoObj.SendingDate,
		SentDate:              dtoObj.SentDate,
		AozoraTransferApplyNo: dtoObj.AozoraTransferApplyNo,
		ApprovalID:            dtoObj.ApprovalID,
		UserID:                dtoObj.UserID,
		User:                  dtoObj.User.ToUserModel(),
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
			DeletedAt: dtoObj.DeletedAt,
		},
	}
}

func ToPayoutDTOs(payouts []*modelPayout.Payout) []*dto.PayoutDTO {
	if payouts == nil {
		return nil
	}

	result := make([]*dto.PayoutDTO, len(payouts))
	for i, payout := range payouts {
		result[i] = ToPayoutDTO(payout)
	}
	return result
}

func ToPayoutModels(dtos []*dto.PayoutDTO) []*modelPayout.Payout {
	if dtos == nil {
		return nil
	}

	result := make([]*modelPayout.Payout, len(dtos))
	for i, dto := range dtos {
		result[i] = ToPayoutModel(dto)
	}
	return result
}
