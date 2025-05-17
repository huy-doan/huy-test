package convert

import (
	modelApprovalStage "github.com/huydq/test/internal/domain/model/approval_stage"
	objectApprovalStage "github.com/huydq/test/internal/domain/object/approval_stage"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	"github.com/huydq/test/internal/infrastructure/persistence/approval_stage/dto"
)

// ToApprovalStageDTO converts an ApprovalStage domain model to an ApprovalStage
func ToApprovalStageDTO(stage *modelApprovalStage.ApprovalStage) *dto.ApprovalStage {
	if stage == nil {
		return nil
	}

	return &dto.ApprovalStage{
		ID:                      stage.ID,
		ApprovalID:              stage.ApprovalID,
		ApprovalWorkflowStageID: stage.ApprovalWorkflowStageID,
		ApproverID:              stage.ApproverID,
		ApprovalResult:          int(stage.ApprovalResult),
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: stage.CreatedAt,
			UpdatedAt: stage.UpdatedAt,
			DeletedAt: stage.DeletedAt,
		},
	}
}

// ToApprovalStageModel converts an ApprovalStage to an ApprovalStage domain model
func ToApprovalStageModel(dtoObj *dto.ApprovalStage) *modelApprovalStage.ApprovalStage {
	if dtoObj == nil {
		return nil
	}

	return &modelApprovalStage.ApprovalStage{
		ID:                      dtoObj.ID,
		ApprovalID:              dtoObj.ApprovalID,
		ApprovalWorkflowStageID: dtoObj.ApprovalWorkflowStageID,
		ApproverID:              dtoObj.ApproverID,
		ApprovalResult:          objectApprovalStage.ApprovalResult(dtoObj.ApprovalResult),
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
			DeletedAt: dtoObj.DeletedAt,
		},
	}
}

// ToApprovalStageDTOs converts a list of ApprovalStage domain models to a list of ApprovalStageDTOs
func ToApprovalStageDTOs(stages []*modelApprovalStage.ApprovalStage) []*dto.ApprovalStage {
	if stages == nil {
		return nil
	}

	result := make([]*dto.ApprovalStage, len(stages))
	for i, stage := range stages {
		result[i] = ToApprovalStageDTO(stage)
	}
	return result
}

// ToApprovalStageModels converts a list of ApprovalStageDTOs to a list of ApprovalStage domain models
func ToApprovalStageModels(dtos []*dto.ApprovalStage) []*modelApprovalStage.ApprovalStage {
	if dtos == nil {
		return nil
	}

	result := make([]*modelApprovalStage.ApprovalStage, len(dtos))
	for i, dto := range dtos {
		result[i] = ToApprovalStageModel(dto)
	}
	return result
}
