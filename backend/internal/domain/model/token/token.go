package token

import (
	"time"

	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

type Token struct {
	ID int

	Token     string
	IsActive  bool
	ExpiredAt time.Time

	util.BaseColumnTimestamp
}

func (t *Token) Invalidate() {
	t.IsActive = false
}
