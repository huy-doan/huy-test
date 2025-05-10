package mapper

import (
	"errors"
	"strconv"
	"time"

	"github.com/huydq/test/internal/datastructure/inputdata"
	"github.com/huydq/test/internal/datastructure/outputdata"
	"github.com/labstack/echo/v4"
)

// parseDateFromString parses a date string in various formats including Japanese format
func parseDateFromString(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC1123,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006年01月02日",          // Japanese format: YYYY年MM月DD日
		"2006年1月2日",            // Japanese format without leading zeros
		"2006年01月02日 15:04:05", // Japanese format with time
	}

	for _, format := range formats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, errors.New("invalid date format")
}

// ToListAuditLogInput converts an Echo Context to a ListAuditLogInput
func ToListAuditLogInput(ctx echo.Context) (*inputdata.ListAuditLogInput, error) {
	input := &inputdata.ListAuditLogInput{}

	// Extract pagination and sorting parameters
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	input.Page = page

	pageSize, err := strconv.Atoi(ctx.QueryParam("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	input.PageSize = pageSize

	input.SortField = ctx.QueryParam("sort_field")
	input.SortOrder = ctx.QueryParam("sort_order")

	// Extract filter parameters
	createdAtStr := ctx.QueryParam("created_at")
	if createdAtStr != "" {
		createdAt, err := parseDateFromString(createdAtStr)
		if err != nil {
			return nil, err
		}
		input.CreatedAt = &createdAt
	}

	userIDStr := ctx.QueryParam("user_id")
	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return nil, err
		}
		input.UserID = &userID
	}

	description := ctx.QueryParam("description")
	if description != "" {
		input.Description = &description
	}

	auditLogType := ctx.QueryParam("audit_log_type")
	if auditLogType != "" {
		input.AuditLogType = &auditLogType
	}

	return input, nil
}

func MapAuditLogListOutputToResponse(output *outputdata.ListAuditLogOutput) map[string]interface{} {
	return map[string]interface{}{
		"audit_logs":  MapAuditLogsToResponse(output.AuditLogs),
		"page":        output.Page,
		"page_size":   output.PageSize,
		"total_pages": output.TotalPages,
		"total":       output.Total,
	}
}

// MapAuditLogsToResponse converts AuditLogOutput slice to a response format
func MapAuditLogsToResponse(auditLogs []*outputdata.AuditLogOutput) []map[string]interface{} {
	if auditLogs == nil {
		return nil
	}

	response := make([]map[string]interface{}, len(auditLogs))
	for i, auditLog := range auditLogs {
		response[i] = MapAuditLogToResponse(auditLog)
	}
	return response
}

// MapAuditLogToResponse converts a single AuditLogOutput to a response map
func MapAuditLogToResponse(auditLog *outputdata.AuditLogOutput) map[string]interface{} {
	if auditLog == nil {
		return nil
	}

	response := map[string]interface{}{
		"id":             auditLog.ID,
		"user_id":        auditLog.UserID,
		"audit_log_type": auditLog.AuditLogType,
		"description":    auditLog.Description,
		"created_at":     auditLog.CreatedAt.Format("2006年01月02日 15:04:05"),
		"updated_at":     auditLog.UpdatedAt.Format("2006年01月02日 15:04:05"),
	}

	if auditLog.TransactionID != nil {
		response["transaction_id"] = *auditLog.TransactionID
	}

	if auditLog.PayoutID != nil {
		response["payout_id"] = *auditLog.PayoutID
	}

	if auditLog.PayinID != nil {
		response["payin_id"] = *auditLog.PayinID
	}

	if auditLog.UserAgent != nil {
		response["user_agent"] = *auditLog.UserAgent
	}

	if auditLog.IPAddress != nil {
		response["ip_address"] = *auditLog.IPAddress
	}

	return response
}

// MapUsersToResponse converts a slice of User models to a response map
// func MapUsersToResponse(users []*user.User) map[string]interface{} {
// 	if users == nil {
// 		return map[string]interface{}{
// 			"users": []interface{}{},
// 		}
// 	}

// 	usersResponse := make([]map[string]interface{}, len(users))
// 	for i, user := range users {
// 		usersResponse[i] = MapUserToResponse(user)
// 	}

// 	return map[string]interface{}{
// 		"users": usersResponse,
// 	}
// }

// MapUserToResponse converts a single User model to a response map
// func MapUserToResponse(user *user.User) map[string]interface{} {
// 	if user == nil {
// 		return nil
// 	}

// 	return map[string]interface{}{
// 		"id":         user.ID,
// 		"first_name": user.FirstName,
// 		"last_name":  user.LastName,
// 		"email":      user.Email,
// 		"full_name":  user.GetFullName(),
// 	}
// }
