package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseMonthYear parses month-year string input
func ParseMonthYear(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, "-")

	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid date format, expected MM-YYYY")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month")
	}
	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year")
	}

	if month < 1 || month > 12 {
		return time.Time{}, fmt.Errorf("month must be between 01 and 12")
	}

	if year <= 1900 || year >= 2500 {
		return time.Time{}, fmt.Errorf("year out of range")
	}

	// Normalize to month start
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}

// FormatMonthYear formats date to month-year
func FormatMonthYear(t time.Time) string {
	// Normalize date to month start
	t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	m := (int)(t.Month())

	return fmt.Sprintf("%02d-%04d", m, t.Year())
}
