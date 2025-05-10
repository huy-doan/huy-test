package base_model

import "time"

type BaseColumnTimestamp struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
