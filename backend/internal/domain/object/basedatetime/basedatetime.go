package object

import (
	"time"
)

/** ----------------------------------------------------------
 * Base Column Timestamp
 * ---------------------------------------------------------- */
type BaseColumnTimestamp struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
