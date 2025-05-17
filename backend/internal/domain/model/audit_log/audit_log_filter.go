package model

import (
	"fmt"
	"time"

	"github.com/huydq/test/internal/domain/model/util"
)

type AuditLogFilter struct {
	util.BaseFilter
	UserID       *int
	AuditLogType *string
	CreatedAt    *time.Time
	Description  *string
}

func NewAuditLogFilter() *AuditLogFilter {
	filter := &AuditLogFilter{}
	filter.ValidSortFields = map[string]bool{
		"id":          true,
		"created_at":  true,
		"description": true,
		"user_id":     true,
		"ip_address":  true,
		"user_agent":  true,
	}

	filter.SetPagination(1, 10)

	return filter
}

func (f *AuditLogFilter) ApplyFilters() {
	f.AddDateFilter("created_at", util.Equal, f.CreatedAt)
	if f.UserID != nil {
		f.AddCondition("user_id", util.Equal, f.UserID)
	}
	if f.AuditLogType != nil {
		f.AddCondition("audit_log_type", util.Equal, f.AuditLogType)
	}
	if f.Description != nil {
		f.AddCondition("description", util.Like, fmt.Sprintf("%%%s%%", *f.Description))
	}
}
