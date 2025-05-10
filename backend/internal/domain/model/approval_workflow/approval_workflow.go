package approval_workflow

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

// ApprovalWorkflowParams contains parameters for creating a new ApprovalWorkflow
type ApprovalWorkflowParams struct {
	ID int
	util.BaseColumnTimestamp
	Name string
}

// ApprovalWorkflow represents an approval workflow definition
type ApprovalWorkflow struct {
	ID int
	util.BaseColumnTimestamp

	Name string
}

// NewApprovalWorkflow creates a new approval workflow
func NewApprovalWorkflow(params ApprovalWorkflowParams) *ApprovalWorkflow {
	return &ApprovalWorkflow{
		ID:                  params.ID,
		Name:                params.Name,
		BaseColumnTimestamp: params.BaseColumnTimestamp,
	}
}
