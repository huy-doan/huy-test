package object

type ApprovalResult int

const (
	ApprovalResultApproved ApprovalResult = 1 // 承認
	ApprovalResultRejected ApprovalResult = 2 // 却下
)

// IsApproved checks if the stage is approved
func (a ApprovalResult) IsApproved() bool {
	return a == ApprovalResultApproved
}

// IsRejected checks if the stage is rejected
func (a ApprovalResult) IsRejected() bool {
	return a == ApprovalResultRejected
}
