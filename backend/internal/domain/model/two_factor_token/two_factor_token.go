package two_factor_token

import (
	"time"

	"github.com/huydq/test/internal/domain/model/user"
	util "github.com/huydq/test/internal/domain/object/basedatetime"
)

type TwoFactorToken struct {
	ID        int
	UserID    int
	Token     string
	MFAType   int
	User      *user.User
	IsUsed    bool
	ExpiredAt time.Time

	util.BaseColumnTimestamp
}
