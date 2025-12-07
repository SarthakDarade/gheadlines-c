package handlers

import (
	"gheadlines/config"
	"gheadlines/db"
	"gheadlines/models"
	"html/template"
	"net/http"
	"time"
)

// UserProfileHandler serves the user profile page
func UserProfileHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUser(r, dbClient, cfg)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		cookie, _ := r.Cookie("access_token")
		profile, err := dbClient.GetProfile(r.Context(), user.ID, cookie.Value)
		if err != nil {
			// Fallback if profile missing
			profile = &models.Profile{
				ID:        user.ID,
				FullName:  user.FullName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		data := struct {
			Profile        models.Profile
			Stats          map[string]int
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
			Profile: *profile,
			Stats: map[string]int{
				"ArticlesRead": 142,
				"Bookmarks":    28,
				"Comments":     15,
			},
			Title:          "Profile - " + cfg.SiteName,
			Description:    "User profile",
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

		tmpl, err := template.ParseFiles(
			"web/templates/user/profile.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// UserEditProfileHandler serves the edit profile page
func UserEditProfileHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUser(r, dbClient, cfg)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		cookie, _ := r.Cookie("access_token")
		profile, err := dbClient.GetProfile(r.Context(), user.ID, cookie.Value)
		if err != nil {
			profile = &models.Profile{
				ID:        user.ID,
				FullName:  user.FullName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		data := struct {
			Profile        models.Profile
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
			Profile:        *profile,
			Title:          "Edit Profile - " + cfg.SiteName,
			Description:    "Edit your profile",
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

		tmpl, err := template.ParseFiles(
			"web/templates/user/edit_profile.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// UserDashboardHandler serves the user dashboard
func UserDashboardHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUser(r, dbClient, cfg)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		cookie, _ := r.Cookie("access_token")
		profile, err := dbClient.GetProfile(r.Context(), user.ID, cookie.Value)
		if err != nil {
			profile = &models.Profile{
				ID:        user.ID,
				FullName:  user.FullName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		data := struct {
			Profile        models.Profile
			Stats          map[string]int
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
			Profile: *profile,
			Stats: map[string]int{
				"ArticlesRead": 142,
				"Bookmarks":    28,
				"Comments":     15,
			},
			Title:          "Dashboard - " + cfg.SiteName,
			Description:    "User dashboard",
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

		tmpl, err := template.ParseFiles(
			"web/templates/user/dashboard.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// UserSettingsHandler serves the account settings page
func UserSettingsHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetCurrentUser(r, dbClient, cfg)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		cookie, _ := r.Cookie("access_token")
		profile, err := dbClient.GetProfile(r.Context(), user.ID, cookie.Value)
		if err != nil {
			profile = &models.Profile{
				ID:        user.ID,
				FullName:  user.FullName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
		data := struct {
			Profile        models.Profile
			Email          string
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
			Profile:        *profile,
			Email:          user.Email,
			Title:          "Settings - " + cfg.SiteName,
			Description:    "Account settings",
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

		tmpl, err := template.ParseFiles(
			"web/templates/user/settings.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// SearchHandler serves the search page
func SearchHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := GetCurrentUser(r, dbClient, cfg)

		query := r.URL.Query().Get("q")
		var articles []models.Article
		if query != "" {
			var err error
			articles, err = dbClient.SearchArticles(r.Context(), query)
			if err != nil {
				// Handle error (maybe log it) but show page
				articles = []models.Article{}
			}
		}

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
			Articles       []models.Article
			Query          string
			CurrentYear    int
		}{
			Title:          "Search results for \"" + query + "\" - " + cfg.SiteName,
			Description:    "Search results",
			SiteName:       cfg.SiteName,
			SiteURL:        cfg.SiteURL,
			CurrentPath:    r.URL.Path,
			CurrentDate:    time.Now().Format("Monday, January 2, 2006"),
			ActiveCategory: "",
			ImageURL:       "",
			JSONLD:         "",
			User:           user,
			Articles:       articles,
			Query:          query,
			CurrentYear:    time.Now().Year(),
		}

		if query == "" {
			data.Title = "Search - " + cfg.SiteName
		}

		tmpl, err := template.ParseFiles(
			"web/templates/search.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// SubscriptionHandler serves the subscription page
func SubscriptionHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _ := GetCurrentUser(r, dbClient, cfg)
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
			Title:          "Subscribe - " + cfg.SiteName,
			Description:    "Premium subscription plans",
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

		tmpl, err := template.ParseFiles(
			"web/templates/subscription.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

// UserUpdateProfileHandler handles profile updates
func UserUpdateProfileHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user, err := GetCurrentUser(r, dbClient, cfg)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		profile := &models.Profile{
			ID:         user.ID,
			FullName:   r.FormValue("fullname"),
			Username:   r.FormValue("username"),
			Bio:        r.FormValue("bio"),
			Occupation: r.FormValue("occupation"),
			Location:   r.FormValue("location"),
			Website:    r.FormValue("website"),
			AvatarURL:  r.FormValue("avatar_url"),
		}

		cookie, _ := r.Cookie("access_token")
		if err := dbClient.UpdateUserProfile(r.Context(), profile, cookie.Value); err != nil {
			http.Redirect(w, r, "/user/edit?error=update_failed", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/user/profile?success=true", http.StatusSeeOther)
	}
}
