package utils

import (
	"encoding/xml"
	"fmt"
	"gheadlines/models"
	"time"
)

// JSONLDArticle represents Schema.org Article
type JSONLDArticle struct {
	Context          string   `json:"@context"`
	Type             string   `json:"@type"`
	MainEntityOfPage string   `json:"mainEntityOfPage"`
	Headline         string   `json:"headline"`
	Image            []string `json:"image"`
	DatePublished    string   `json:"datePublished"`
	DateModified     string   `json:"dateModified"`
	Author           struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		URL  string `json:"url,omitempty"`
	} `json:"author"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo struct {
			Type string `json:"@type"`
			URL  string `json:"url"`
		} `json:"logo"`
	} `json:"publisher"`
	Description    string `json:"description"`
	ArticleSection string `json:"articleSection,omitempty"`
	Keywords       string `json:"keywords,omitempty"`
}

// GenerateJSONLD generates JSON-LD for an article
func GenerateJSONLD(article *models.Article, siteURL string, siteName string) JSONLDArticle {
	return JSONLDArticle{
		Context:          "https://schema.org",
		Type:             "NewsArticle",
		MainEntityOfPage: fmt.Sprintf("%s/article/%s", siteURL, article.Slug),
		Headline:         article.Title,
		Image:            []string{article.ImageURL},
		DatePublished:    article.CreatedAt.Format(time.RFC3339),
		DateModified:     article.UpdatedAt.Format(time.RFC3339),
		Author: struct {
			Type string `json:"@type"`
			Name string `json:"name"`
			URL  string `json:"url,omitempty"`
		}{
			Type: "Person",
			Name: article.Author,
		},
		Publisher: struct {
			Type string `json:"@type"`
			Name string `json:"name"`
			Logo struct {
				Type string `json:"@type"`
				URL  string `json:"url"`
			} `json:"logo"`
		}{
			Type: "Organization",
			Name: siteName,
			Logo: struct {
				Type string `json:"@type"`
				URL  string `json:"url"`
			}{
				Type: "ImageObject",
				URL:  siteURL + "/static/img/logo.png", // Ensure this exists or use text
			},
		},
		Description:    article.Excerpt,
		ArticleSection: article.Category,
		Keywords:       article.Category + ", News, World News",
	}
}

// SitemapURL represents a URL in sitemap
type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

// URLSet represents the sitemap root
type URLSet struct {
	XMLName xmlns        `xml:"urlset"`
	URLs    []SitemapURL `xml:"url"`
}

type xmlns struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
}

// GenerateSitemap generates XML sitemap
func GenerateSitemap(articles []models.Article, siteURL string) ([]byte, error) {
	urlSet := URLSet{
		URLs: []SitemapURL{
			{
				Loc:        siteURL + "/",
				LastMod:    time.Now().Format("2006-01-02"),
				ChangeFreq: "hourly",
				Priority:   "1.0",
			},
		},
	}

	for _, article := range articles {
		urlSet.URLs = append(urlSet.URLs, SitemapURL{
			Loc:        fmt.Sprintf("%s/article/%s", siteURL, article.Slug),
			LastMod:    article.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "daily",
			Priority:   "0.8",
		})
	}

	return xml.MarshalIndent(urlSet, "", "  ")
}
