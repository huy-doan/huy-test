package object

import (
	"time"
)

/** ----------------------------------------------------------
 * Base Column Timestamp
 * ---------------------------------------------------------- */
type BaseColumnTimestamp struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
