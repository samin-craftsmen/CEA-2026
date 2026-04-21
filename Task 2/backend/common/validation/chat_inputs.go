package validation

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var mealTypePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,31}$`)

func UserID(value, fieldName string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s is required", fieldName)
	}
	if strings.ContainsAny(trimmed, " \t\r\n") {
		return "", fmt.Errorf("%s must not contain whitespace", fieldName)
	}
	return trimmed, nil
}

func Date(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("date is required")
	}
	parsed, err := time.Parse("2006-01-02", trimmed)
	if err != nil || parsed.Format("2006-01-02") != trimmed {
		return "", fmt.Errorf("invalid date %q: expected YYYY-MM-DD", value)
	}
	return trimmed, nil
}

func MealType(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "", fmt.Errorf("meal_type is required")
	}
	if !mealTypePattern.MatchString(normalized) {
		return "", fmt.Errorf("invalid meal_type %q: use 1-32 lowercase letters, numbers, hyphens, or underscores", value)
	}
	return normalized, nil
}

func Status(value string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	if normalized == "" {
		return "", fmt.Errorf("status is required")
	}
	if normalized != "YES" && normalized != "NO" {
		return "", fmt.Errorf("invalid status %q: must be YES or NO", value)
	}
	return normalized, nil
}

func Location(value string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	if normalized == "" {
		return "", fmt.Errorf("location is required")
	}
	if normalized != "OFFICE" && normalized != "WFH" {
		return "", fmt.Errorf("invalid location %q: must be OFFICE or WFH", value)
	}
	return normalized, nil
}

func DayStatusType(value string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	if normalized == "" {
		return "", fmt.Errorf("type is required")
	}
	switch normalized {
	case "GOVERNMENT_HOLIDAY", "OFFICE_CLOSED", "SPECIAL_EVENT":
		return normalized, nil
	default:
		return "", fmt.Errorf("invalid type %q: must be GOVERNMENT_HOLIDAY, OFFICE_CLOSED, or SPECIAL_EVENT", value)
	}
}

func Note(value string, required bool) (string, error) {
	trimmed := strings.TrimSpace(value)
	if required && trimmed == "" {
		return "", fmt.Errorf("note is required for SPECIAL_EVENT")
	}
	return trimmed, nil
}
