package models

const (
	MFATypeEmail int = 1 // Email-based MFA - メール
)

func GetMFATypeTitle(mfaType int) string {
	switch mfaType {
	case MFATypeEmail:
		return "Email"
	default:
		return "Email"
	}
}
