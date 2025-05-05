package messages

const (
	// Error Code
	CodeValidationError      = "VALIDATION_ERROR"
	CodeNotFound             = "NOT_FOUND"
	CodeUnauthorized         = "UNAUTHORIZED"
	CodeForbidden            = "FORBIDDEN"
	CodeBadRequest           = "BAD_REQUEST"
	CodeInternalError        = "INTERNAL_ERROR"
	CodeDatabaseError        = "DATABASE_ERROR"
	CodeExternalServiceError = "EXTERNAL_SERVICE_ERROR"
	CodeUnknownError         = "UNKNOWN_ERROR"
	CodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	CodeTimeout              = "TIMEOUT"
	CodeTooManyRequests      = "TOO_MANY_REQUESTS"
	CodeUnsupportedMediaType = "UNSUPPORTED_MEDIA_TYPE"
	CodeMethodNotAllowed     = "METHOD_NOT_ALLOWED"
	CodeNotAcceptable        = "NOT_ACCEPTABLE"

	// Error Types
	TypeValidationError      = "VALIDATION"
	TypeAuthorizationError   = "AUTHORIZATION"
	TypeDatabaseError        = "DATABASE"
	TypeExternalServiceError = "EXTERNAL"
	TypeServerError          = "SERVER"
	TypeClientError          = "CLIENT"
	TypeNotFoundError        = "NOT_FOUND"
	TypeUnknownError         = "UNKNOWN"
	TypeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	TypeTimeoutError         = "TIMEOUT"

	// Error Messages
	MsgInvalidCredentials = "メールアドレスまたはパスワードが無効です"
	MsgEmailAlreadyExists = "このメールアドレスは既に登録されています"
	MsgUserNotFound       = "ユーザーが見つかりません"
	MsgInvalidPassword    = "現在のパスワードが正しくありません"
	MsgAccountLocked      = "アカウントがロックされています"
)
