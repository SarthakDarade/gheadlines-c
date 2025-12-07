package handlers

import (
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

// AdminDashboardHandler renders the admin dashboard
func AdminDashboardHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch articles (limit 50 for dashboard)
		articles, err := dbClient.GetArticles(r.Context(), 50, 0, nil, "")
		if err != nil {
			// Handle error or use empty list
			articles = []models.Article{} // Fallback
		}

		tmpl, err := template.ParseFiles(
			"web/templates/admin/dashboard.html",
			"web/templates/partials/admin_header.html",  // Assuming we'll create this
			"web/templates/partials/admin_sidebar.html", // Assuming we'll create this
		)
		if err != nil {
			// Try parsing just the dashboard if partials fail (during dev)
			tmpl, err = template.ParseFiles("web/templates/admin/dashboard.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		data := struct {
			SiteName  string
			Articles  []models.Article
			User      string
			CSRFToken string
		}{
			SiteName:  cfg.SiteName,
			Articles:  articles,
			User:      cfg.AdminUsername,
			CSRFToken: r.Context().Value(CSRFTokenKey).(string),
		}

		tmpl.Execute(w, data)
	}
}

// AdminNewArticleHandler renders the create article form
func AdminNewArticleHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, _ := dbClient.GetCategories(r.Context(), "")

		tmpl, err := template.ParseFiles("web/templates/admin/new.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Categories []models.Category
			CSRFToken  string
		}{
			Categories: categories,
			CSRFToken:  r.Context().Value(CSRFTokenKey).(string),
		}

		tmpl.Execute(w, data)
	}
}

// AdminEditArticleHandler renders the edit article form
func AdminEditArticleHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/adm/edit/")
		if id == "" {
			http.Redirect(w, r, "/adm/dashboard", http.StatusSeeOther)
			return
		}

		article, err := dbClient.GetArticleByID(r.Context(), id, "")
		if err != nil {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}

		categories, _ := dbClient.GetCategories(r.Context(), "")

		tmpl, err := template.ParseFiles("web/templates/admin/edit.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Article    *models.Article
			Categories []models.Category
			CSRFToken  string
		}{
			Article:    article,
			Categories: categories,
			CSRFToken:  r.Context().Value(CSRFTokenKey).(string),
		}

		tmpl.Execute(w, data)
	}
}

// AdminSaveArticleHandler handles creating or updating an article
func AdminSaveArticleHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		id := r.FormValue("id")
		title := r.FormValue("title")
		content := r.FormValue("content")
		excerpt := r.FormValue("excerpt")
		category := r.FormValue("category")
		imageURL := r.FormValue("image_url")
		// status := r.FormValue("status") // draft or published - unused for now

		// Auto-generate slug if new or empty
		slug := r.FormValue("slug")
		if slug == "" {
			slug = utils.GenerateSlug(title)
		}

		// Create article object
		article := models.Article{
			ID:        id,
			Title:     title,
			Content:   content,
			Excerpt:   excerpt,
			Category:  category,
			ImageURL:  imageURL,
			Slug:      slug,
			UpdatedAt: time.Now(),
			// Status: status, // Need to add Status to Article model if not present
		}

		if id == "" {
			// Create new
			article.CreatedAt = time.Now()
			article.Date = time.Now().Format("2006-01-02")

			if err := dbClient.CreateArticle(r.Context(), &article); err != nil {
				log.Printf("Error creating article: %v", err)
				http.Error(w, "Failed to create article", http.StatusInternalServerError)
				return
			}
		} else {
			// Update existing
			// Try to preserve existing fields like CreatedAt, Views, etc.
			existing, err := dbClient.GetArticleByID(r.Context(), id, "")
			if err == nil && existing != nil {
				article.CreatedAt = existing.CreatedAt
				article.Date = existing.Date
				article.Views = existing.Views
				article.Likes = existing.Likes
			}

			if err := dbClient.UpdateArticle(r.Context(), &article); err != nil {
				log.Printf("Error updating article: %v", err)
				http.Error(w, "Failed to update article", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/adm/dashboard", http.StatusSeeOther)
	}
}
