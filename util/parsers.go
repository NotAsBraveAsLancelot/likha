package util

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseCount converts a user-friendly string like "10k", "10m", "10b" to an int64.
// It is case-insensitive and trims whitespace.
func ParseCount(s string) (int64, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	var multiplier int64 = 1

	if strings.HasSuffix(s, "k") {
		multiplier = 1_000
		s = strings.TrimSuffix(s, "k")
	} else if strings.HasSuffix(s, "m") {
		multiplier = 1_000_000
		s = strings.TrimSuffix(s, "m")
	} else if strings.HasSuffix(s, "b") {
		multiplier = 1_000_000_000
		s = strings.TrimSuffix(s, "b")
	}

	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: '%s'", s)
	}

	return val * multiplier, nil
}

// InterfaceToInt safely converts an interface{} to int.
// It's useful for parsing settings from the YAML config.
func InterfaceToInt(v interface{}) (int, bool) {
	// YAML unmarshals numbers as int, float64, etc. Handle both.
	switch i := v.(type) {
	case int:
		return i, true
	case float64:
		return int(i), true // Note: this truncates
	default:
		return 0, false
	}
}

// InterfaceToInt64 safely converts an interface{} to int64.
func InterfaceToInt64(v interface{}) (int64, bool) {
	switch i := v.(type) {
	case int:
		return int64(i), true
	case int64:
		return i, true
	case float64:
		return int64(i), true
	default:
		return 0, false
	}
}

// InterfaceToFloat64 safely converts an interface{} to float64.
func InterfaceToFloat64(v interface{}) (float64, bool) {
	switch f := v.(type) {
	case float64:
		return f, true
	case int:
		return float64(f), true
	case int64:
		return float64(f), true
	default:
		return 0, false
	}
}
