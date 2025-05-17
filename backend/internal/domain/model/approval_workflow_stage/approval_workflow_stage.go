package approval_workflow_stage

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// ApprovalWorkflowStageParams contains parameters for creating a new ApprovalWorkflowStage
type ApprovalWorkflowStageParams struct {
	ID int
	util.BaseColumnTimestamp
	WorkflowID     int
	StageName      string
	Level          int
	ApproverRoleID int
	ApproverCount  int
}

// ApprovalWorkflowStage represents a stage in an approval workflow
type ApprovalWorkflowStage struct {
	ID int
	util.BaseColumnTimestamp

	WorkflowID     int
	StageName      string
	Level          int
	ApproverRoleID int
	ApproverCount  int
}

// NewApprovalWorkflowStage creates a new approval workflow stage
func NewApprovalWorkflowStage(params ApprovalWorkflowStageParams) *ApprovalWorkflowStage {
	return &ApprovalWorkflowStage{
		ID:                  params.ID,
		WorkflowID:          params.WorkflowID,
		StageName:           params.StageName,
		Level:               params.Level,
		ApproverRoleID:      params.ApproverRoleID,
		ApproverCount:       params.ApproverCount,
		BaseColumnTimestamp: params.BaseColumnTimestamp,
	}
}
