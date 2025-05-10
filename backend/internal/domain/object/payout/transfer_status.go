package object

// TransferStatus represents the status of a transfer
type TransferStatus int

const (
	TransferStatusInProgress     TransferStatus = 1 // 振込中
	TransferStatusWhitelistError TransferStatus = 2 // ホワイトリスト追加エラー
	TransferStatusApiError       TransferStatus = 3 // 振込依頼APIエラー
	TransferStatusFailed         TransferStatus = 4 // 振込依頼失敗
	TransferStatusRequested      TransferStatus = 5 // 振込依頼済み
	TransferStatusProcessed      TransferStatus = 6 // 送金手続き済み
)

// String returns the string representation of the transfer status
func (t TransferStatus) String() string {
	switch t {
	case TransferStatusInProgress:
		return "振込中"
	case TransferStatusWhitelistError:
		return "ホワイトリスト追加エラー"
	case TransferStatusApiError:
		return "振込依頼APIエラー"
	case TransferStatusFailed:
		return "振込依頼失敗"
	case TransferStatusRequested:
		return "振込依頼済み"
	case TransferStatusProcessed:
		return "送金手続き済み"
	default:
		return "不明"
	}
}

// IsError checks if the transfer status is an error
func (t TransferStatus) IsError() bool {
	return t == TransferStatusWhitelistError ||
		t == TransferStatusApiError ||
		t == TransferStatusFailed
}

// IsProcessed checks if the transfer is processed
func (t TransferStatus) IsProcessed() bool {
	return t == TransferStatusProcessed
}

// IsRequested checks if the transfer is requested
func (t TransferStatus) IsRequested() bool {
	return t == TransferStatusRequested
}

// IsInProgress checks if the transfer is in progress
func (t TransferStatus) IsInProgress() bool {
	return t == TransferStatusInProgress
}
