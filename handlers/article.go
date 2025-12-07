package handlers

import (
	"encoding/json"
	"gheadlines/config"
	"gheadlines/db"
	"gheadlines/internal/utils"
	"gheadlines/models"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

// ArticleData holds data for the article page template
type ArticleData struct {
	Article         *models.Article
	RelatedArticles []models.Article
	Categories      []models.Category
	MetaTags        template.HTML
	JSONLD          template.HTML
	Breadcrumbs     template.HTML
	CanonicalURL    string
	CurrentYear     int
	User            *models.User
}

// ArticleHandler handles individual article pages
// ArticleHandler handles individual article pages
func ArticleHandler(dbClient *db.Client, siteURL string, siteName string, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract slug from URL path
		slug := strings.TrimPrefix(r.URL.Path, "/article/")
		if slug == "" {
			http.NotFound(w, r)
			return
		}

		// Get access token from cookie
		var accessToken string
		if cookie, err := r.Cookie("access_token"); err == nil {
			accessToken = cookie.Value
		}

		// Fetch article
		article, err := dbClient.GetArticleBySlug(r.Context(), slug, accessToken)
		if err != nil || article == nil {
			log.Printf("Article not found for slug: %s, error: %v", slug, err)
			http.NotFound(w, r)
			return
		}
		// Convert content to HTML
		article.ContentHTML = template.HTML(article.Content)

		// Fetch related articles (same category, excluding current)
		relatedArticles := []models.Article{}
		if article.Category != "" {
			allArticles, _ := dbClient.GetArticles(r.Context(), 10, 0, &article.Category, accessToken)
			for _, a := range allArticles {
				if a.ID != article.ID && len(relatedArticles) < 3 {
					relatedArticles = append(relatedArticles, a)
				}
			}
		}

		// Fetch categories for navigation
		categories, _ := dbClient.GetCategories(r.Context(), accessToken)

		// Generate JSON-LD
		jsonLD := utils.GenerateJSONLD(article, siteURL, siteName)
		jsonLDBytes, _ := json.Marshal(jsonLD)

		// Get current user
		user, _ := GetCurrentUser(r, dbClient, cfg)

		data := struct {
			Article         *models.Article
			RelatedArticles []models.Article
			Categories      []models.Category
			SiteName        string
			SiteURL         string
			CurrentDate     string
			JSONLD          template.HTML
			Title           string
			Description     string
			ImageURL        string
			CurrentPath     string
			ActiveCategory  string
			User            *models.User
			CurrentYear     int
		}{
			Article:         article,
			RelatedArticles: relatedArticles,
			Categories:      categories,
			SiteName:        siteName,
			SiteURL:         siteURL,
			CurrentDate:     time.Now().Format("Monday, January 2, 2006"),
			JSONLD:          template.HTML(jsonLDBytes),
			Title:           article.Title,
			Description:     article.Excerpt,
			ImageURL:        article.ImageURL,
			CurrentPath:     r.URL.Path,
			ActiveCategory:  article.CategoryObj.Slug,
			User:            user,
			CurrentYear:     time.Now().Year(),
		}

		// Render template
		tmpl, err := template.ParseFiles(
			"web/templates/article.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data)
	}
}
