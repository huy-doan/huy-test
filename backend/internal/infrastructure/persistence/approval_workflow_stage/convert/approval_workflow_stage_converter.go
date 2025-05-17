package convert

import (
	modelApprovalWorkflowStage "github.com/huydq/test/internal/domain/model/approval_workflow_stage"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	"github.com/huydq/test/internal/infrastructure/persistence/approval_workflow_stage/dto"
)

// ToApprovalWorkflowStageDTO converts an ApprovalWorkflowStage domain model to an ApprovalWorkflowStage
func ToApprovalWorkflowStageDTO(stage *modelApprovalWorkflowStage.ApprovalWorkflowStage) *dto.ApprovalWorkflowStage {
	if stage == nil {
		return nil
	}

	return &dto.ApprovalWorkflowStage{
		ID:             stage.ID,
		WorkflowID:     stage.WorkflowID,
		StageName:      stage.StageName,
		Level:          stage.Level,
		ApproverRoleID: stage.ApproverRoleID,
		ApproverCount:  stage.ApproverCount,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: stage.CreatedAt,
			UpdatedAt: stage.UpdatedAt,
			DeletedAt: stage.DeletedAt,
		},
	}
}

// ToApprovalWorkflowStageModel converts an ApprovalWorkflowStage to an ApprovalWorkflowStage domain model
func ToApprovalWorkflowStageModel(dtoObj *dto.ApprovalWorkflowStage) *modelApprovalWorkflowStage.ApprovalWorkflowStage {
	if dtoObj == nil {
		return nil
	}

	return &modelApprovalWorkflowStage.ApprovalWorkflowStage{
		ID:             dtoObj.ID,
		WorkflowID:     dtoObj.WorkflowID,
		StageName:      dtoObj.StageName,
		Level:          dtoObj.Level,
		ApproverRoleID: dtoObj.ApproverRoleID,
		ApproverCount:  dtoObj.ApproverCount,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
			DeletedAt: dtoObj.DeletedAt,
		},
	}
}

// ToApprovalWorkflowStageDTOs converts a list of ApprovalWorkflowStage domain models to a list of ApprovalWorkflowStageDTOs
func ToApprovalWorkflowStageDTOs(stages []*modelApprovalWorkflowStage.ApprovalWorkflowStage) []*dto.ApprovalWorkflowStage {
	if stages == nil {
		return nil
	}

	result := make([]*dto.ApprovalWorkflowStage, len(stages))
	for i, stage := range stages {
		result[i] = ToApprovalWorkflowStageDTO(stage)
	}
	return result
}

// ToApprovalWorkflowStageModels converts a list of ApprovalWorkflowStageDTOs to a list of ApprovalWorkflowStage domain models
func ToApprovalWorkflowStageModels(dtos []*dto.ApprovalWorkflowStage) []*modelApprovalWorkflowStage.ApprovalWorkflowStage {
	if dtos == nil {
		return nil
	}

	result := make([]*modelApprovalWorkflowStage.ApprovalWorkflowStage, len(dtos))
	for i, dto := range dtos {
		result[i] = ToApprovalWorkflowStageModel(dto)
	}
	return result
}
