package models

const (
	ApprovalResultApproved int = 1 // 承認
	ApprovalResultRejected int = 2 // 却下
)

type ApprovalStage struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	ApprovalID              int `json:"approval_id"`
	ApprovalWorkflowStageID int `json:"approval_workflow_stage_id"`
	ApproverID              int `json:"approver_id"`
	ApprovalResult          int `json:"approval_result"`

	Approval              *Approval              `json:"approval,omitempty" gorm:"foreignKey:ApprovalID"`
	ApprovalWorkflowStage *ApprovalWorkflowStage `json:"approval_workflow_stage,omitempty" gorm:"foreignKey:ApprovalWorkflowStageID"`
	Approver              *User                  `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
}

func (ApprovalStage) TableName() string {
	return "approval_stage"
}
