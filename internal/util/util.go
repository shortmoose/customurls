package util

import (
	"regexp"
	"strings"
)

// GetKey returns a lowercase version of the given string, returns an empty string if invalid.
func GetKey(str string) string {
	if !regexp.MustCompile(`^[A-Za-z0-9-]+$`).MatchString(str) {
		return ""
	}

	return strings.ToLower(str)
}
