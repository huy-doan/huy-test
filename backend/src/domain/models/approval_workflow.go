package models

type ApprovalWorkflow struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	Name string `json:"name"`

	ApprovalWorkflowStages []ApprovalWorkflowStage `json:"approval_workflow_stages,omitempty" gorm:"foreignKey:WorkflowID"`
	Approvals              []Approval              `json:"approvals,omitempty" gorm:"foreignKey:ApprovalWorkflowID"`
}

func (ApprovalWorkflow) TableName() string {
	return "approval_workflow"
}
