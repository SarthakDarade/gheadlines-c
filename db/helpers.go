package db

import (
	"gheadlines/models"
	"strings"
	"unicode"
)

// generateSlug creates a URL-friendly slug from title and ID
func generateSlug(title, id string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	var result strings.Builder
	for _, r := range slug {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result.WriteRune(r)
		}
	}

	slug = result.String()
	// Remove multiple dashes
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Add ID suffix for uniqueness
	if id != "" {
		slug += "-" + id[:8]
	}

	return slug
}

// getCategoryByName returns a Category object for a category name
func getCategoryByName(categoryName string) *models.Category {
	categoryMap := map[string]models.Category{
		"Technology":    {Name: "Technology", Slug: "technology", Color: "#3B82F6"},
		"Environment":   {Name: "Environment", Slug: "environment", Color: "#10B981"},
		"Business":      {Name: "Business", Slug: "business", Color: "#F59E0B"},
		"Health":        {Name: "Health", Slug: "health", Color: "#EF4444"},
		"Sports":        {Name: "Sports", Slug: "sports", Color: "#8B5CF6"},
		"Culture":       {Name: "Culture", Slug: "culture", Color: "#EC4899"},
		"Politics":      {Name: "Politics", Slug: "politics", Color: "#DC2626"},
		"Science":       {Name: "Science", Slug: "science", Color: "#7C3AED"},
		"Entertainment": {Name: "Entertainment", Slug: "entertainment", Color: "#EC4899"},
	}

	if cat, ok := categoryMap[categoryName]; ok {
		return &cat
	}

	// Default category
	return &models.Category{
		Name:  categoryName,
		Slug:  strings.ToLower(strings.ReplaceAll(categoryName, " ", "-")),
		Color: "#6B7280",
	}
}
