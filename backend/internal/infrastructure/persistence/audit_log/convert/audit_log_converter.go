package convert

import (
	model "github.com/huydq/test/internal/domain/model/audit_log"
	object "github.com/huydq/test/internal/domain/object/audit_log"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	"github.com/huydq/test/internal/infrastructure/persistence/audit_log/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
	"gorm.io/gorm"
)

func ToAuditLogDTO(auditLog *model.AuditLog) *dto.AuditLogDTO {
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

	result := &dto.AuditLogDTO{
		ID:            auditLog.ID,
		UserID:        auditLog.UserID,
		AuditLogType:  string(auditLog.AuditLogType),
		Description:   auditLog.Description,
		TransactionID: auditLog.TransactionID,
		PayoutID:      auditLog.PayoutID,
		PayinID:       auditLog.PayinID,
		UserAgent:     userAgent,
		IPAddress:     ipAddress,
		BaseColumnTimestamp: persistence.BaseColumnTimestamp{
			CreatedAt: auditLog.CreatedAt,
			UpdatedAt: auditLog.UpdatedAt,
		},
	}

	// Handle the conversion from *time.Time to gorm.DeletedAt
	if auditLog.DeletedAt != nil {
		result.DeletedAt = gorm.DeletedAt{
			Time:  *auditLog.DeletedAt,
			Valid: true,
		}
	} else {
		result.DeletedAt = gorm.DeletedAt{
			Valid: false,
		}
	}

	return result
}

func ToAuditLogModel(dtoObj *dto.AuditLogDTO) *model.AuditLog {
	if dtoObj == nil {
		return nil
	}

	var userAgent *object.UserAgent
	var ipAddress *object.IPAddress

	if dtoObj.UserAgent != nil {
		ua := object.UserAgent(*dtoObj.UserAgent)
		userAgent = &ua
	}

	if dtoObj.IPAddress != nil {
		ip := object.IPAddress(*dtoObj.IPAddress)
		ipAddress = &ip
	}

	// Create base domain model
	result := &model.AuditLog{
		ID:            dtoObj.ID,
		UserID:        dtoObj.UserID,
		AuditLogType:  object.AuditLogType(dtoObj.AuditLogType),
		Description:   dtoObj.Description,
		TransactionID: dtoObj.TransactionID,
		PayoutID:      dtoObj.PayoutID,
		PayinID:       dtoObj.PayinID,
		UserAgent:     userAgent,
		IPAddress:     ipAddress,
		BaseColumnTimestamp: util.BaseColumnTimestamp{
			CreatedAt: dtoObj.CreatedAt,
			UpdatedAt: dtoObj.UpdatedAt,
		},
	}

	// Handle the conversion from gorm.DeletedAt to *time.Time
	if dtoObj.DeletedAt.Valid {
		deletedAt := dtoObj.DeletedAt.Time
		result.DeletedAt = &deletedAt
	} else {
		result.DeletedAt = nil
	}

	return result
}

func ToAuditLogDTOs(auditLogs []*model.AuditLog) []*dto.AuditLogDTO {
	result := make([]*dto.AuditLogDTO, len(auditLogs))
	for i, auditLog := range auditLogs {
		result[i] = ToAuditLogDTO(auditLog)
	}
	return result
}

func ToAuditLogModels(dtos []*dto.AuditLogDTO) []*model.AuditLog {
	result := make([]*model.AuditLog, len(dtos))
	for i, dto := range dtos {
		result[i] = ToAuditLogModel(dto)
	}
	return result
}
