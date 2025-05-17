package object

// BankAccountType represents the type of bank account
type BankAccountType int

const (
	BankAccountTypeOrdinary BankAccountType = 1 // 普通預金
	BankAccountTypeCurrent  BankAccountType = 2 // 当座預金
	BankAccountTypeFixed    BankAccountType = 3 // 定期預金
)

func (b BankAccountType) String() string {
	switch b {
	case BankAccountTypeOrdinary:
		return "普通預金"
	case BankAccountTypeCurrent:
		return "当座預金"
	case BankAccountTypeFixed:
		return "定期預金"
	default:
		return "不明"
	}
}

func (b BankAccountType) IsValid() bool {
	return b == BankAccountTypeOrdinary ||
		b == BankAccountTypeCurrent ||
		b == BankAccountTypeFixed
}
