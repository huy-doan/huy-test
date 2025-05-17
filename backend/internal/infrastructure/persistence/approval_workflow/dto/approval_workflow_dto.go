package dto

import (
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

type ApprovalWorkflow struct {
	ID int `json:"id"`
	util.BaseColumnTimestamp

	Name string `json:"name"`
}

func (ApprovalWorkflow) TableName() string {
	return "approval_workflow"
}
