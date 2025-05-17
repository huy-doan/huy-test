package dto

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	userDto "github.com/huydq/test/internal/infrastructure/persistence/user/dto"
)

const (
	ApprovalResultApproved int = 1 // 承認
	ApprovalResultRejected int = 2 // 却下
)

type ApprovalStage struct {
	ID int `json:"id"`
	util.BaseColumnTimestamp

	ApprovalID              int `json:"approval_id"`
	ApprovalWorkflowStageID int `json:"approval_workflow_stage_id"`
	ApproverID              int `json:"approver_id"`
	ApprovalResult          int `json:"approval_result"`

	Approver *userDto.User `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
}

func (ApprovalStage) TableName() string {
	return "approval_stage"
}
