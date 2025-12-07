package services

import (
	"encoding/json"
	"fmt"
	"gheadlines/models"
	"html/template"
	"strings"
	"time"
)

// MetaTags holds all SEO meta information
type MetaTags struct {
	Title       string
	Description string
	Keywords    string
	ImageURL    string
	URL         string
	Type        string
	SiteName    string
	PublishedAt time.Time
	UpdatedAt   time.Time
	Author      string
}

// GenerateJSONLD creates JSON-LD structured data for an article
func GenerateJSONLD(article *models.Article, siteURL string, siteName string) template.HTML {
	articleURL := fmt.Sprintf("%s/article/%s", siteURL, article.Slug)
	
	imageURL := article.ImageURL
	if imageURL == "" {
		imageURL = fmt.Sprintf("%s/static/og-default.jpg", siteURL)
	}

	authorName := article.Author
	if authorName == "" {
		authorName = "Global Headlines"
	}

	jsonLD := map[string]interface{}{
		"@context": "https://schema.org",
		"@type":    "NewsArticle",
		"headline": article.Title,
		"image": []string{
			imageURL,
		},
		"datePublished": article.CreatedAt.Format(time.RFC3339),
		"dateModified":  article.UpdatedAt.Format(time.RFC3339),
		"author": map[string]interface{}{
			"@type": "Person",
			"name":  authorName,
		},
		"publisher": map[string]interface{}{
			"@type": "Organization",
			"name":  siteName,
			"logo": map[string]interface{}{
				"@type": "ImageObject",
				"url":   fmt.Sprintf("%s/static/logo.png", siteURL),
			},
		},
		"mainEntityOfPage": map[string]interface{}{
			"@type": "WebPage",
			"@id":   articleURL,
		},
		"articleBody": article.Content,
		"description": article.Excerpt,
	}

	jsonBytes, _ := json.MarshalIndent(jsonLD, "", "  ")
	return template.HTML(fmt.Sprintf("<script type=\"application/ld+json\">\n%s\n</script>", string(jsonBytes)))
}

// GenerateBreadcrumbs creates BreadcrumbList JSON-LD
func GenerateBreadcrumbs(items []BreadcrumbItem, siteURL string) template.HTML {
	breadcrumbList := []map[string]interface{}{}
	
	for i, item := range items {
		breadcrumbList = append(breadcrumbList, map[string]interface{}{
			"@type":    "ListItem",
			"position": i + 1,
			"name":     item.Name,
			"item":     fmt.Sprintf("%s%s", siteURL, item.URL),
		})
	}

	jsonLD := map[string]interface{}{
		"@context":        "https://schema.org",
		"@type":           "BreadcrumbList",
		"itemListElement": breadcrumbList,
	}

	jsonBytes, _ := json.MarshalIndent(jsonLD, "", "  ")
	return template.HTML(fmt.Sprintf("<script type=\"application/ld+json\">\n%s\n</script>", string(jsonBytes)))
}

// BreadcrumbItem represents a breadcrumb navigation item
type BreadcrumbItem struct {
	Name string
	URL  string
}

// GenerateMetaTags creates meta tag HTML
func GenerateMetaTags(tags MetaTags) template.HTML {
	var meta strings.Builder

	// Basic meta tags
	meta.WriteString(fmt.Sprintf("<title>%s</title>\n", template.HTMLEscapeString(tags.Title)))
	meta.WriteString(fmt.Sprintf("<meta name=\"description\" content=\"%s\">\n", template.HTMLEscapeString(tags.Description)))
	if tags.Keywords != "" {
		meta.WriteString(fmt.Sprintf("<meta name=\"keywords\" content=\"%s\">\n", template.HTMLEscapeString(tags.Keywords)))
	}

	// Canonical URL
	meta.WriteString(fmt.Sprintf("<link rel=\"canonical\" href=\"%s\">\n", template.HTMLEscapeString(tags.URL)))

	// Open Graph tags
	meta.WriteString(fmt.Sprintf("<meta property=\"og:type\" content=\"%s\">\n", template.HTMLEscapeString(tags.Type)))
	meta.WriteString(fmt.Sprintf("<meta property=\"og:title\" content=\"%s\">\n", template.HTMLEscapeString(tags.Title)))
	meta.WriteString(fmt.Sprintf("<meta property=\"og:description\" content=\"%s\">\n", template.HTMLEscapeString(tags.Description)))
	meta.WriteString(fmt.Sprintf("<meta property=\"og:url\" content=\"%s\">\n", template.HTMLEscapeString(tags.URL)))
	meta.WriteString(fmt.Sprintf("<meta property=\"og:image\" content=\"%s\">\n", template.HTMLEscapeString(tags.ImageURL)))
	meta.WriteString(fmt.Sprintf("<meta property=\"og:site_name\" content=\"%s\">\n", template.HTMLEscapeString(tags.SiteName)))
	if !tags.PublishedAt.IsZero() {
		meta.WriteString(fmt.Sprintf("<meta property=\"article:published_time\" content=\"%s\">\n", tags.PublishedAt.Format(time.RFC3339)))
	}
	if !tags.UpdatedAt.IsZero() {
		meta.WriteString(fmt.Sprintf("<meta property=\"article:modified_time\" content=\"%s\">\n", tags.UpdatedAt.Format(time.RFC3339)))
	}
	if tags.Author != "" {
		meta.WriteString(fmt.Sprintf("<meta property=\"article:author\" content=\"%s\">\n", template.HTMLEscapeString(tags.Author)))
	}

	// Twitter Card tags
	meta.WriteString("<meta name=\"twitter:card\" content=\"summary_large_image\">\n")
	meta.WriteString(fmt.Sprintf("<meta name=\"twitter:title\" content=\"%s\">\n", template.HTMLEscapeString(tags.Title)))
	meta.WriteString(fmt.Sprintf("<meta name=\"twitter:description\" content=\"%s\">\n", template.HTMLEscapeString(tags.Description)))
	meta.WriteString(fmt.Sprintf("<meta name=\"twitter:image\" content=\"%s\">\n", template.HTMLEscapeString(tags.ImageURL)))

	return template.HTML(meta.String())
}

// GenerateSitemapURL creates a sitemap URL entry
type SitemapURL struct {
	Loc        string
	LastMod    time.Time
	ChangeFreq string
	Priority   float64
}

// FormatSitemapURL formats a URL for sitemap.xml
func FormatSitemapURL(url SitemapURL) string {
	return fmt.Sprintf(`  <url>
    <loc>%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>%s</changefreq>
    <priority>%.1f</priority>
  </url>`, url.Loc, url.LastMod.Format("2006-01-02"), url.ChangeFreq, url.Priority)
}


