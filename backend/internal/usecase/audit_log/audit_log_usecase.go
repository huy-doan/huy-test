package usecase

import (
	"context"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	model "github.com/huydq/test/internal/domain/model/audit_log"
	"github.com/huydq/test/internal/domain/model/util"
	object "github.com/huydq/test/internal/domain/object/audit_log"
	"github.com/huydq/test/internal/domain/service"
)

type AuditLogUsecase interface {
	// Create creates a new audit log
	Create(ctx context.Context, input *inputdata.CreateAuditLogInput) error

	// List retrieves a paginated list of audit logs based on filter criteria
	List(ctx context.Context, input *inputdata.ListAuditLogInput) (*outputdata.ListAuditLogOutput, error)
}

type auditLogUsecaseImpl struct {
	auditLogService service.AuditLogService
}

func NewAuditLogUsecase(auditLogService service.AuditLogService) AuditLogUsecase {
	return &auditLogUsecaseImpl{
		auditLogService: auditLogService,
	}
}

// Create handles the creation of a new audit log
func (uc *auditLogUsecaseImpl) Create(ctx context.Context, input *inputdata.CreateAuditLogInput) error {
	auditLogModel := uc.createInputToModel(input)
	return uc.auditLogService.Create(ctx, auditLogModel)
}

// List handles retrieving a filtered, paginated list of audit logs
func (uc *auditLogUsecaseImpl) List(ctx context.Context, input *inputdata.ListAuditLogInput) (*outputdata.ListAuditLogOutput, error) {
	filter := uc.listInputToFilter(input)

	auditLogs, totalPages, total, err := uc.auditLogService.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	auditLogOutputs := uc.modelsToOutputs(auditLogs)

	return &outputdata.ListAuditLogOutput{
		AuditLogs:  auditLogOutputs,
		Page:       filter.Pagination.Page,
		PageSize:   filter.Pagination.PageSize,
		TotalPages: totalPages,
		Total:      total,
	}, nil
}

// createInputToModel converts a CreateAuditLogInput to an AuditLog domain model
func (uc *auditLogUsecaseImpl) createInputToModel(input *inputdata.CreateAuditLogInput) *model.AuditLog {
	if input == nil {
		return nil
	}

	var userAgent *object.UserAgent
	var ipAddress *object.IPAddress

	if input.UserAgent != nil {
		ua := object.UserAgent(*input.UserAgent)
		userAgent = &ua
	}

	if input.IPAddress != nil {
		ip := object.IPAddress(*input.IPAddress)
		ipAddress = &ip
	}

	return model.NewAuditLog(model.NewAuditLogParams{
		UserID:        input.UserID,
		AuditLogType:  object.AuditLogType(input.AuditLogType),
		Description:   input.Description,
		TransactionID: input.TransactionID,
		PayoutID:      input.PayoutID,
		PayinID:       input.PayinID,
		UserAgent:     userAgent,
		IPAddress:     ipAddress,
	})
}

// listInputToFilter converts a ListAuditLogInput to an AuditLogFilter domain model
func (uc *auditLogUsecaseImpl) listInputToFilter(input *inputdata.ListAuditLogInput) *model.AuditLogFilter {
	if input == nil {
		return model.NewAuditLogFilter()
	}

	filter := model.NewAuditLogFilter()
	filter.UserID = input.UserID
	filter.AuditLogType = input.AuditLogType
	filter.CreatedAt = input.CreatedAt
	filter.Description = input.Description

	// Setting pagination
	filter.SetPagination(input.Page, input.PageSize)

	// Setting sorting
	if input.SortField != "" {
		sortOrder := util.MapSortDirection(input.SortOrder)
		filter.SetSort(input.SortField, sortOrder)
	}

	return filter
}

// modelToOutput converts an AuditLog domain model to an AuditLogOutput
func (uc *auditLogUsecaseImpl) modelToOutput(auditLog *model.AuditLog) *outputdata.AuditLogOutput {
	if auditLog == nil {
		return nil
	}

	var userAgent, ipAddress *string
	if auditLog.UserAgent != nil {
		ua := auditLog.UserAgent.String()
		userAgent = &ua
	}
	if auditLog.IPAddress != nil {
		ip := auditLog.IPAddress.String()
		ipAddress = &ip
	}

	return &outputdata.AuditLogOutput{
		ID:            auditLog.ID,
		UserID:        auditLog.UserID,
		AuditLogType:  string(auditLog.AuditLogType),
		Description:   auditLog.Description,
		TransactionID: auditLog.TransactionID,
		PayoutID:      auditLog.PayoutID,
		PayinID:       auditLog.PayinID,
		UserAgent:     userAgent,
		IPAddress:     ipAddress,
		CreatedAt:     auditLog.CreatedAt,
		UpdatedAt:     auditLog.UpdatedAt,
	}
}

// modelsToOutputs converts a slice of AuditLog domain models to a slice of AuditLogOutputs
func (uc *auditLogUsecaseImpl) modelsToOutputs(auditLogs []*model.AuditLog) []*outputdata.AuditLogOutput {
	if auditLogs == nil {
		return nil
	}

	outputs := make([]*outputdata.AuditLogOutput, len(auditLogs))
	for i, auditLog := range auditLogs {
		outputs[i] = uc.modelToOutput(auditLog)
	}
	return outputs
}
