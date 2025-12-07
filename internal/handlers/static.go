package handlers

import (
	"gheadlines/config"
	"gheadlines/db"
	rootHandlers "gheadlines/handlers"
	"gheadlines/models"
	"html/template"
	"net/http"
	"time"
)

// StaticPageHandler serves static templates like About, Contact, Privacy, Terms
func StaticPageHandler(pageName string, dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current user
		user, _ := rootHandlers.GetCurrentUser(r, dbClient, cfg)

		// Data for the template
		data := struct {
			Title          string
			Description    string
			SiteName       string
			SiteURL        string
			CurrentPath    string
			CurrentDate    string
			ActiveCategory string
			ImageURL       string
			JSONLD         template.HTML
			User           *models.User
			CurrentYear    int
		}{
			SiteName:       cfg.SiteName,
			SiteURL:        cfg.SiteURL,
			CurrentPath:    r.URL.Path,
			CurrentDate:    time.Now().Format("Monday, January 2, 2006"),
			ActiveCategory: "",
			ImageURL:       "",
			JSONLD:         "",
			User:           user,
			CurrentYear:    time.Now().Year(),
		}

		// Set specific titles based on page
		switch pageName {
		case "about":
			data.Title = "About Us"
			data.Description = "Learn about Global Headlines, our mission, and our team."
		case "contact":
			data.Title = "Contact Us"
			data.Description = "Get in touch with Global Headlines."
		case "privacy":
			data.Title = "Privacy Policy"
			data.Description = "Our commitment to protecting your privacy."
		case "terms":
			data.Title = "Terms of Service"
			data.Description = "Rules and regulations for using Global Headlines."
		}

		// Parse templates
		tmpl, err := template.ParseFiles(
			"web/templates/static/"+pageName+".html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute template
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
