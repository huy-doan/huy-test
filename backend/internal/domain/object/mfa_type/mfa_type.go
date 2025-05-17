package object

type MFAType int

const (
	MFA_TYPE_EMAIL MFAType = 1
)

// String returns the string representation of the MFAType
func (m MFAType) String() string {
	switch m {
	case MFA_TYPE_EMAIL:
		return "Email"
	default:
		return "Unknown"
	}
}
