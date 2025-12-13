package handlers

import (
	"encoding/xml"
	"fmt"
	"gheadlines/db"
	"gheadlines/models"
	"net/http"
	"time"
)

// Sitemap represents the sitemap.xml structure
type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// NewsSitemap represents the Google News sitemap structure
type NewsSitemap struct {
	XMLName   xml.Name  `xml:"urlset"`
	Xmlns     string    `xml:"xmlns,attr"`
	XmlnsNews string    `xml:"xmlns:news,attr"`
	URLs      []NewsURL `xml:"url"`
}

// NewsURL represents a URL entry in the news sitemap
type NewsURL struct {
	Loc  string   `xml:"loc"`
	News NewsMeta `xml:"news:news"`
}

// NewsMeta represents the news metadata
type NewsMeta struct {
	Publication NewsPublication `xml:"news:publication"`
	PubDate     string          `xml:"news:publication_date"`
	Title       string          `xml:"news:title"`
}

// NewsPublication represents the publisher info
type NewsPublication struct {
	Name     string `xml:"news:name"`
	Language string `xml:"news:language"`
}

// URL represents a URL entry in the sitemap
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}

// SitemapHandler generates and serves sitemap.xml
func SitemapHandler(dbClient *db.Client, siteURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Fetch all articles
		articles, err := dbClient.GetArticles(r.Context(), 1000, 0, nil, "")
		if err != nil {
			http.Error(w, "Failed to generate sitemap", http.StatusInternalServerError)
			return
		}

		// Build sitemap
		sitemap := Sitemap{
			Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
			URLs:  []URL{},
		}

		// Add homepage
		sitemap.URLs = append(sitemap.URLs, URL{
			Loc:        siteURL + "/",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "daily",
			Priority:   1.0,
		})

		// Add article pages
		for _, article := range articles {
			lastMod := article.UpdatedAt
			if lastMod.IsZero() {
				lastMod = article.CreatedAt
			}

			sitemap.URLs = append(sitemap.URLs, URL{
				Loc:        fmt.Sprintf("%s/article/%s", siteURL, article.Slug),
				LastMod:    lastMod.Format("2006-01-02"),
				ChangeFreq: "weekly",
				Priority:   0.8,
			})
		}

		// Fetch categories and add category pages
		categories, err := dbClient.GetCategories(r.Context(), "")
		if err == nil {
			for _, category := range categories {
				sitemap.URLs = append(sitemap.URLs, URL{
					Loc:        fmt.Sprintf("%s/category/%s", siteURL, category.Slug),
					LastMod:    time.Now().Format("2006-01-02"),
					ChangeFreq: "daily",
					Priority:   0.7,
				})
			}
		}

		// Set XML content type
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		// Write XML
		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "  ")
		if err := encoder.Encode(sitemap); err != nil {
			http.Error(w, "Failed to encode sitemap", http.StatusInternalServerError)
			return
		}
	}
}

// NewsSitemapHandler generates and serves sitemap-news.xml for Google News
func NewsSitemapHandler(dbClient *db.Client, siteURL string, siteName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch articles from last 48 hours for Google News
		// Note: dbClient.GetArticles might need a time filter, but for now we fetch latest and filter
		// In a real optimized app, you'd pass a "since" parameter to the DB query
		limit := 100
		articles, err := dbClient.GetArticles(r.Context(), limit, 0, nil, "")
		if err != nil {
			http.Error(w, "Failed to generate news sitemap", http.StatusInternalServerError)
			return
		}

		// Filter for last 2 days
		cutoff := time.Now().Add(-48 * time.Hour)
		var recentArticles []models.Article
		for _, a := range articles {
			// Check CreatedAt or UpdatedAt
			date := a.CreatedAt
			if !a.UpdatedAt.IsZero() {
				date = a.UpdatedAt
			}
			if date.After(cutoff) {
				recentArticles = append(recentArticles, a)
			}
		}

		// Build sitemap
		sitemap := NewsSitemap{
			Xmlns:     "http://www.sitemaps.org/schemas/sitemap/0.9",
			XmlnsNews: "http://www.google.com/schemas/sitemap-news/0.9",
			URLs:      []NewsURL{},
		}

		for _, article := range recentArticles {
			pubDate := article.CreatedAt
			if !pubDate.IsZero() {
				loc := fmt.Sprintf("%s/article/%s", siteURL, article.Slug)
				if article.Slug == "" {
					loc = fmt.Sprintf("%s/article/%s", siteURL, article.ID)
				}
				sitemap.URLs = append(sitemap.URLs, NewsURL{
					Loc: loc,
					News: NewsMeta{
						Publication: NewsPublication{
							Name:     siteName,
							Language: "en",
						},
						PubDate: pubDate.Format("2006-01-02"), // YYYY-MM-DD
						Title:   article.Title,
					},
				})
			}
		}

		// Set XML content type
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		// Write XML
		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "  ")
		if err := encoder.Encode(sitemap); err != nil {
			logError(w, "Failed to encode news sitemap", err)
			return
		}
	}
}

func logError(w http.ResponseWriter, msg string, err error) {
	fmt.Printf("%s: %v\n", msg, err)
	// Don't write to w if it's already written
}
