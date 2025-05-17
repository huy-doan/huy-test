package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONField is a generic type for handling JSON fields in database
type JSONField[T any] map[string]T

// Value implements the driver.Valuer interface for database serialization
func (jf JSONField[T]) Value() (driver.Value, error) {
	return json.Marshal(jf)
}

// Scan implements the sql.Scanner interface for database deserialization
func (jf *JSONField[T]) Scan(value any) error {
	if value == nil {
		*jf = make(JSONField[T])
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, jf)
}
