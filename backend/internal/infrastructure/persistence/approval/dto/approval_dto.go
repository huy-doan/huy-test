package dto

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	approvalStageDto "github.com/huydq/test/internal/infrastructure/persistence/approval_stage/dto"
	approvalWorkflowDto "github.com/huydq/test/internal/infrastructure/persistence/approval_workflow/dto"
)

const (
	ApprovalStatusPending      int = 1 // 承認待ち
	ApprovalStatusWaitApproval int = 2 // 承認中
	ApprovalStatusApproved     int = 3 // 承認済み
	ApprovalStatusRejected     int = 4 // 却下
)

type Approval struct {
	ID int `json:"id"`
	util.BaseColumnTimestamp

	ApprovalWorkflowID int `json:"approval_workflow_id"`
	ApprovalStatus     int `json:"approval_status"`

	ApprovalWorkflow *approvalWorkflowDto.ApprovalWorkflow `json:"approval_workflow,omitempty" gorm:"foreignKey:ApprovalWorkflowID"`
	ApprovalStages   []approvalStageDto.ApprovalStage      `json:"approval_stages,omitempty" gorm:"foreignKey:ApprovalID"`
}

func (Approval) TableName() string {
	return "approval"
}
