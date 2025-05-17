package convert

import (
	modelApprovalWorkflow "github.com/huydq/test/internal/domain/model/approval_workflow"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	"github.com/huydq/test/internal/infrastructure/persistence/approval_workflow/dto"
)

// ToApprovalWorkflowDTO converts an ApprovalWorkflow domain model to an ApprovalWorkflow
func ToApprovalWorkflowDTO(workflow *modelApprovalWorkflow.ApprovalWorkflow) *dto.ApprovalWorkflow {
	if workflow == nil {
		return nil
	}

	return &dto.ApprovalWorkflow{
		ID:   workflow.ID,
		Name: workflow.Name,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: workflow.CreatedAt,
			UpdatedAt: workflow.UpdatedAt,
			DeletedAt: workflow.DeletedAt,
		},
	}
}

// ToApprovalWorkflowModel converts an ApprovalWorkflow to an ApprovalWorkflow domain model
func ToApprovalWorkflowModel(dtoObj *dto.ApprovalWorkflow) *modelApprovalWorkflow.ApprovalWorkflow {
	if dtoObj == nil {
		return nil
	}

	return &modelApprovalWorkflow.ApprovalWorkflow{
		ID:   dtoObj.ID,
		Name: dtoObj.Name,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
			DeletedAt: dtoObj.DeletedAt,
		},
	}
}

// ToApprovalWorkflowDTOs converts a list of ApprovalWorkflow domain models to a list of ApprovalWorkflowDTOs
func ToApprovalWorkflowDTOs(workflows []*modelApprovalWorkflow.ApprovalWorkflow) []*dto.ApprovalWorkflow {
	if workflows == nil {
		return nil
	}

	result := make([]*dto.ApprovalWorkflow, len(workflows))
	for i, workflow := range workflows {
		result[i] = ToApprovalWorkflowDTO(workflow)
	}
	return result
}

// ToApprovalWorkflowModels converts a list of ApprovalWorkflowDTOs to a list of ApprovalWorkflow domain models
func ToApprovalWorkflowModels(dtos []*dto.ApprovalWorkflow) []*modelApprovalWorkflow.ApprovalWorkflow {
	if dtos == nil {
		return nil
	}

	result := make([]*modelApprovalWorkflow.ApprovalWorkflow, len(dtos))
	for i, dto := range dtos {
		result[i] = ToApprovalWorkflowModel(dto)
	}
	return result
}
