package dto

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

type ApprovalWorkflowStage struct {
	ID int `json:"id"`
	util.BaseColumnTimestamp

	WorkflowID     int    `json:"workflow_id"`
	StageName      string `json:"stage_name"`
	Level          int    `json:"level"`
	ApproverRoleID int    `json:"approver_role_id"`
	ApproverCount  int    `json:"approver_count" gorm:"default:1"`
}

func (ApprovalWorkflowStage) TableName() string {
	return "approval_workflow_stage"
}
