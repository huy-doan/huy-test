package messages

const (
	MsgListAuditLogsSuccess    = "監査ログの一覧を取得しました"
	MsgGetAuditLogUsersSuccess = "ユーザー一覧を取得しました"

	// Role Success Messages
	MsgListRolesSuccess                  = "ロール一覧を取得しました"
	MsgGetRoleSuccess                    = "ロールを取得しました"
	MsgCreateRoleSuccess                 = "ロールが作成されました"
	MsgUpdateRoleSuccess                 = "ロールが更新されました"
	MsgDeleteRoleSuccess                 = "ロールが削除されました"
	MsgBatchUpdateRolePermissionsSuccess = "ロール権限の一括更新が完了しました"

	// Permission Success Messages
	MsgListPermissionsSuccess = "権限一覧を取得しました"

	// Auth related success messages
	MsgLoginSuccess        = "ログインしました"
	MsgLogoutSuccess       = "ログアウトしました"
	MsgMFARequired         = "2FA認証が必要です"
	MsgAuthCodeSentSuccess = "認証コードが正常に送信されました"

	// merchant related success messages
	MsgListMerchantsSuccess = "加盟店一覧を取得しました"

	// payout related success messages
	MsgListPayoutsSuccess = "出金履歴の取得に成功しました"

	// User related success messages
	MsgListUsersSuccess  = "ユーザー一覧を取得しました"
	MsgCreateUserSuccess = "ユーザーを登録しました"
	MsgUpdateUserSuccess = "ユーザーを更新しました"
	MsgDeleteUserSuccess = "ユーザーを削除しました"
	MsgGetUserSuccess    = "ユーザーを取得しました"
)
