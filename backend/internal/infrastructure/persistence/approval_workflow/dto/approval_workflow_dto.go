package dto

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

type ApprovalWorkflowDTO struct {
	ID int `json:"id"`
	util.BaseColumnTimestamp

	Name string `json:"name"`
}

func (ApprovalWorkflowDTO) TableName() string {
	return "approval_workflow"
}
