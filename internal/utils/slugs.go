package utils

import (
	"regexp"
	"strings"
)

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove non-alphanumeric characters (except hyphens)
	reg, _ := regexp.Compile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	// Remove multiple hyphens
	reg2, _ := regexp.Compile("-+")
	s = reg2.ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	return s
}
