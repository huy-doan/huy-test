package object

type PayinFileStatus int

const (
	StatusPending PayinFileStatus = 0
	StatusSuccess PayinFileStatus = 1
	StatusFailed  PayinFileStatus = 2
)
