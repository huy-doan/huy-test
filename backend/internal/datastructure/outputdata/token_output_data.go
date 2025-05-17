package outputdata

import (
	"time"
)

type TokenOutputData struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	IsActive  bool      `json:"is_active"`
	ExpiredAt time.Time `json:"expired_at"`
}
