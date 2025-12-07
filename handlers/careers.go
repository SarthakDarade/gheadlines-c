package handlers

import (
	"gheadlines/config"
	"gheadlines/db"
	"gheadlines/models"
	"html/template"
	"log"
	"net/http"
	"time"
)

// CareersHandler renders the careers page
func CareersHandler(dbClient *db.Client, siteURL string, siteName string, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mock positions logic (in real app, fetch from DB)
		positions := []struct {
			Title, Department, Location, Type string
		}{
			{"Senior Political Correspondent", "Editorial", "Washington, DC", "Full-time"},
			{"Tech Editor", "Editorial", "San Francisco, CA / Remote", "Full-time"},
			{"Data Journalist", "Graphics", "New York, NY", "Full-time"},
			{"Frontend Engineer", "Product", "Remote", "Full-time"},
			{"Social Media Manager", "Marketing", "London, UK", "Contract"},
		}

		user, _ := GetCurrentUser(r, dbClient, cfg)

		data := struct {
			SiteName       string
			SiteURL        string
			CurrentDate    string
			Title          string
			CurrentPath    string
			User           *models.User
			Positions      interface{}
			CurrentYear    int
			Description    string
			ImageURL       string
			JSONLD         string
			ActiveCategory string
		}{
			SiteName:       siteName,
			SiteURL:        siteURL,
			CurrentDate:    time.Now().Format("Monday, January 2, 2006"),
			Title:          "Careers",
			CurrentPath:    r.URL.Path,
			User:           user,
			Positions:      positions,
			CurrentYear:    time.Now().Year(),
			Description:    "Join the Global Headlines team and help shape the future of journalism.",
			ImageURL:       "",
			JSONLD:         "",
			ActiveCategory: "careers",
		}

		tmpl, err := template.ParseFiles(
			"web/templates/careers.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Template Execution Error: %v", err)
		}
	}
}

// SubmitCareerApplicationHandler handles the form submission
func SubmitCareerApplicationHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		message := r.FormValue("message")

		resumeURL := r.FormValue("resume_url")

		// Create application object
		app := &models.CareerApplication{
			Name:      name,
			Email:     email,
			Phone:     phone,
			Message:   message,
			ResumeURL: resumeURL,
		}

		// Save to DB
		err := dbClient.CreateCareerApplication(r.Context(), app)
		if err != nil {
			log.Printf("DB Error: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Return Success (Redirect or JSON)
		// For form submission, we can redirect with a success query param
		http.Redirect(w, r, "/careers?success=true", http.StatusSeeOther)
	}
}
