package utils

import (
	"fmt"
	"time"
)

// Helper function to parse time from various formats
func ParseTime(timeStr string) (time.Time, error) {
	// Try RFC3339 format (ISO 8601)
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return t, nil
	}

	// Try RFC3339Nano format
	if t, err := time.Parse(time.RFC3339Nano, timeStr); err == nil {
		return t, nil
	}

	// Try layout with milliseconds
	if t, err := time.Parse("2006-01-02T15:04:05.000Z", timeStr); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unable to parse time string: %s", timeStr)
}
