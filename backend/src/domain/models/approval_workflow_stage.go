package models

type ApprovalWorkflowStage struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	WorkflowID     int    `json:"workflow_id"`
	StageName      string `json:"stage_name"`
	Level          int    `json:"level"`
	ApproverRoleID int    `json:"approver_role_id"`
	ApproverCount  int    `json:"approver_count" gorm:"default:1"`

	ApprovalWorkflow *ApprovalWorkflow `json:"approval_workflow,omitempty" gorm:"foreignKey:WorkflowID"`
	ApproverRole     *Role             `json:"approver_role,omitempty" gorm:"foreignKey:ApproverRoleID"`
	ApprovalStages   []ApprovalStage   `json:"approval_stages,omitempty" gorm:"foreignKey:ApprovalWorkflowStageID"`
}

func (ApprovalWorkflowStage) TableName() string {
	return "approval_workflow_stage"
}
