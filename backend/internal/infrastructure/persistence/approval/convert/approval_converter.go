package convert

import (
	modelApproval "github.com/huydq/test/internal/domain/model/approval"
	objectApproval "github.com/huydq/test/internal/domain/object/approval"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	"github.com/huydq/test/internal/infrastructure/persistence/approval/dto"
)

// ToApprovalDTO converts an Approval domain model to an ApprovalDTO
func ToApprovalDTO(approval *modelApproval.Approval) *dto.ApprovalDTO {
	if approval == nil {
		return nil
	}

	return &dto.ApprovalDTO{
		ID:                 approval.ID,
		ApprovalWorkflowID: approval.ApprovalWorkflowID,
		ApprovalStatus:     int(approval.ApprovalStatus),
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: approval.CreatedAt,
			UpdatedAt: approval.UpdatedAt,
			DeletedAt: approval.DeletedAt,
		},
	}
}

// ToApprovalModel converts an ApprovalDTO to an Approval domain model
func ToApprovalModel(dtoObj *dto.ApprovalDTO) *modelApproval.Approval {
	if dtoObj == nil {
		return nil
	}

	return &modelApproval.Approval{
		ID:                 dtoObj.ID,
		ApprovalWorkflowID: dtoObj.ApprovalWorkflowID,
		ApprovalStatus:     objectApproval.ApprovalStatus(dtoObj.ApprovalStatus),
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
			DeletedAt: dtoObj.DeletedAt,
		},
	}
}

// ToApprovalDTOs converts a list of Approval domain models to a list of ApprovalDTOs
func ToApprovalDTOs(approvals []*modelApproval.Approval) []*dto.ApprovalDTO {
	if approvals == nil {
		return nil
	}

	result := make([]*dto.ApprovalDTO, len(approvals))
	for i, approval := range approvals {
		result[i] = ToApprovalDTO(approval)
	}
	return result
}

// ToApprovalModels converts a list of ApprovalDTOs to a list of Approval domain models
func ToApprovalModels(dtos []*dto.ApprovalDTO) []*modelApproval.Approval {
	if dtos == nil {
		return nil
	}

	result := make([]*modelApproval.Approval, len(dtos))
	for i, dto := range dtos {
		result[i] = ToApprovalModel(dto)
	}
	return result
}
