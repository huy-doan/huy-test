package object

// PayoutStatus represents the status of a payout
type PayoutStatus int

const (
	PayoutStatusDraft     PayoutStatus = 1 // ドラフト
	PayoutStatusCreated   PayoutStatus = 2 // 振込データ作成済み
	PayoutStatusProcessed PayoutStatus = 3 // 送金手続き済み
)

func (p PayoutStatus) String() string {
	switch p {
	case PayoutStatusDraft:
		return "ドラフト"
	case PayoutStatusCreated:
		return "振込データ作成済み"
	case PayoutStatusProcessed:
		return "送金手続き済み"
	default:
		return "不明"
	}
}

func (p PayoutStatus) IsDraft() bool {
	return p == PayoutStatusDraft
}

func (p PayoutStatus) IsCreated() bool {
	return p == PayoutStatusCreated
}

func (p PayoutStatus) IsProcessed() bool {
	return p == PayoutStatusProcessed
}

func GetPayoutStatusFromString(s string) (PayoutStatus, bool) {
	switch s {
	case "ドラフト":
		return PayoutStatusDraft, true
	case "振込データ作成済み":
		return PayoutStatusCreated, true
	case "送金手続き済み":
		return PayoutStatusProcessed, true
	default:
		return 0, false
	}
}

func GetPayoutStatusFromInt(i int) (PayoutStatus, bool) {
	switch i {
	case 1:
		return PayoutStatusDraft, true
	case 2:
		return PayoutStatusCreated, true
	case 3:
		return PayoutStatusProcessed, true
	default:
		return 0, false
	}
}
