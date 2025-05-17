package messages

const (
	ValidateRequired       = "%s フィールドは必須です"
	ValidateEmail          = "%s フィールドは有効なメールアドレスである必要があります"
	ValidateMin            = "%s フィールドは最低 %s 文字以上である必要があります"
	ValidateMax            = "%s フィールドは最大 %s 文字以内である必要があります"
	ValidateKana           = "%s フィールドはカタカナのみを含める必要があります"
	ValidateField          = "%s フィールドの検証に失敗しました: %s"
	ValidatePasswordPolicy = "%s フィールドは、12 文字以上で、大文字、小文字、数字、記号を含める必要があります。"
	MsgInvalidParameter    = "パラメータが無効です"
	MsgInvalidRequestData  = "リクエストデータが無効です"
	MsgValidationFailed    = "入力値の検証に失敗しました"
	MsgInternalError       = "内部サーバーエラーが発生しました"
)
