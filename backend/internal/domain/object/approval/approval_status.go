package object

type ApprovalStatus int

const (
	ApprovalStatusPending      ApprovalStatus = 1 // 承認待ち
	ApprovalStatusWaitApproval ApprovalStatus = 2 // 承認中
	ApprovalStatusApproved     ApprovalStatus = 3 // 承認済み
	ApprovalStatusRejected     ApprovalStatus = 4 // 却下
)

func (a ApprovalStatus) IsPending() bool {
	return a == ApprovalStatusPending
}

func (a ApprovalStatus) IsInProgress() bool {
	return a == ApprovalStatusWaitApproval
}

func (a ApprovalStatus) IsApproved() bool {
	return a == ApprovalStatusApproved
}

func (a ApprovalStatus) IsRejected() bool {
	return a == ApprovalStatusRejected
}
