package approval_stage

import (
	object "github.com/huydq/test/internal/domain/object/approval_stage"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// ApprovalStageParams contains parameters for creating a new ApprovalStage
type ApprovalStageParams struct {
	ID int
	util.BaseColumnTimestamp
	ApprovalID              int
	ApprovalWorkflowStageID int
	ApproverID              int
	ApprovalResult          object.ApprovalResult
}

// ApprovalStage represents a stage in an approval process
type ApprovalStage struct {
	ID int
	util.BaseColumnTimestamp

	ApprovalID              int
	ApprovalWorkflowStageID int
	ApproverID              int
	ApprovalResult          object.ApprovalResult
}

// NewApprovalStage creates a new approval stage
func NewApprovalStage(params ApprovalStageParams) *ApprovalStage {
	return &ApprovalStage{
		ID:                      params.ID,
		ApprovalID:              params.ApprovalID,
		ApprovalWorkflowStageID: params.ApprovalWorkflowStageID,
		ApproverID:              params.ApproverID,
		ApprovalResult:          params.ApprovalResult,
		BaseColumnTimestamp:     params.BaseColumnTimestamp,
	}
}

func (a *ApprovalStage) Approve() {
	a.ApprovalResult = object.ApprovalResultApproved
}

func (a *ApprovalStage) Reject() {
	a.ApprovalResult = object.ApprovalResultRejected
}
