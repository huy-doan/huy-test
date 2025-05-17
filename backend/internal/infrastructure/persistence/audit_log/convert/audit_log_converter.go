package convert

import (
	model "github.com/huydq/test/internal/domain/model/audit_log"
	userModel "github.com/huydq/test/internal/domain/model/user"
	object "github.com/huydq/test/internal/domain/object/audit_log"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
	"github.com/huydq/test/internal/infrastructure/persistence/audit_log/dto"
	userDto "github.com/huydq/test/internal/infrastructure/persistence/user/dto"
	persistence "github.com/huydq/test/internal/infrastructure/persistence/util"
	"gorm.io/gorm"
)

func ToAuditLogDTO(auditLog *model.AuditLog) *dto.AuditLog {
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

	var user *userDto.User
	if auditLog.User != nil {
		user = userDto.ToUserDTO(auditLog.User)
	}

	result := &dto.AuditLog{
		ID:            auditLog.ID,
		User:          user,
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

func ToAuditLogModel(dtoObj *dto.AuditLog) *model.AuditLog {
	if dtoObj == nil {
		return nil
	}

	var userAgent *object.UserAgent
	var ipAddress *object.IPAddress
	var user *userModel.User

	if dtoObj.UserAgent != nil {
		ua := object.UserAgent(*dtoObj.UserAgent)
		userAgent = &ua
	}

	if dtoObj.IPAddress != nil {
		ip := object.IPAddress(*dtoObj.IPAddress)
		ipAddress = &ip
	}

	if dtoObj.User != nil {
		user = dtoObj.User.ToUserModel()
	}

	// Create base domain model
	result := &model.AuditLog{
		ID:            dtoObj.ID,
		UserID:        dtoObj.UserID,
		User:          user,
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

func ToAuditLogDTOs(auditLogs []*model.AuditLog) []*dto.AuditLog {
	result := make([]*dto.AuditLog, len(auditLogs))
	for i, auditLog := range auditLogs {
		result[i] = ToAuditLogDTO(auditLog)
	}
	return result
}

func ToAuditLogModels(dtos []*dto.AuditLog) []*model.AuditLog {
	result := make([]*model.AuditLog, len(dtos))
	for i, dto := range dtos {
		result[i] = ToAuditLogModel(dto)
	}
	return result
}
