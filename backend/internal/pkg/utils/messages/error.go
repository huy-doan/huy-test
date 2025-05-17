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
	MsgBadRequest         = "不正なリクエストです"

	// Route and API errors
	MsgRouteNotFound    = "リクエストされたURLが見つかりません"
	MsgMethodNotAllowed = "このリクエストメソッドは許可されていません"

	// Role Error Messages
	MsgListRolesError                  = "ロール一覧の取得に失敗しました"
	MsgGetRoleError                    = "ロールの取得に失敗しました"
	MsgRoleNotFoundError               = "ロールが見つかりません"
	MsgCreateRoleError                 = "ロールの作成に失敗しました"
	MsgUpdateRoleError                 = "ロールの更新に失敗しました"
	MsgDeleteRoleError                 = "ロールの削除に失敗しました"
	MsgBatchUpdateRolePermissionsError = "ロール権限の一括更新に失敗しました"
	MsgIDRequiredError                 = "IDが必要です"

	// Permission Error Messages
	MsgListPermissionsError = "権限一覧の取得に失敗しました"

	// Auth middleware error messages
	MsgUnauthorized       = "認証ヘッダーが必要です"
	MsgInvalidToken       = "トークンの検証エラー"
	MsgForbidden          = "このリソースにアクセスする権限がありません"
	MsgTokenBlacklisted   = "無効または失効したトークン"
	MsgInvalidTokenFormat = "認証ヘッダーの形式が正しくありません"
	MsgLogoutFailed       = "ログアウトに失敗しました"

	// Auth related error messages
	MsgLoginFailed      = "ログインに失敗しました"
	MsgUnauthenticated  = "認証されていません"
	MsgMFAVerifyFailed  = "認証コードの検証に失敗しました"
	MsgResendCodeFailed = "認証コードの再送信に失敗しました"

	// user related error messages
	MsgCreateUserFailed = "ユーザーを登録できませんでした"
	MsgUpdateUserFailed = "ユーザーを更新できませんでした"
	MsgDeleteUserFailed = "ユーザーを削除できませんでした"
	MsgGetUserFailed    = "ユーザーを取得できませんでした"
)
