package commands

import (
	"sort"
	"strings"
)

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func normalizeUserRef(value string) string {
	trimmed := strings.TrimSpace(value)
	switch {
	case strings.HasPrefix(trimmed, "<") || strings.HasSuffix(trimmed, ">"):
		if strings.HasPrefix(trimmed, "<users/") && strings.HasSuffix(trimmed, ">") {
			trimmed = strings.TrimSuffix(strings.TrimPrefix(trimmed, "<users/"), ">")
		} else {
			return ""
		}
	case strings.HasPrefix(trimmed, "users/"):
		trimmed = strings.TrimPrefix(trimmed, "users/")
	}

	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" || strings.ContainsAny(trimmed, "<> \t\r\n") {
		return ""
	}
	return trimmed
}
