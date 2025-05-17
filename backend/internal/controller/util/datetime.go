package util

import (
	"errors"
	"time"
)

func ParseDateFromString(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC1123,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006年01月02日",          // Japanese format: YYYY年MM月DD日
		"2006年1月2日",            // Japanese format without leading zeros
		"2006年01月02日 15:04:05", // Japanese format with time
	}

	for _, format := range formats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, errors.New("invalid date format")
}
