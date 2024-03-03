package internal

import (
	"regexp"
	"strings"
)

func CreateSlug(s string) string {
	slug := strings.ReplaceAll(strings.ToLower(s), " ", "-")
	regex := regexp.MustCompile("[^a-zA-Z0-9-]")
	return regex.ReplaceAllString(slug, "")
}
