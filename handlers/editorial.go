package handlers

import (
	"gheadlines/config"
	"gheadlines/db"
	"gheadlines/models"
	"html/template"
	"net/http"
	"strings"
	"time"
)

// EditorialTeamHandler serves the team listing page
func EditorialTeamHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		team, err := dbClient.GetEditorialTeam(r.Context())
		if err != nil {
			team = []models.EditorialTeamMember{}
		}

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
			Team           []models.EditorialTeamMember
		}{
			Title:          "Editorial Team - " + cfg.SiteName,
			Description:    "Meet the team behind " + cfg.SiteName,
			SiteName:       cfg.SiteName,
			SiteURL:        cfg.SiteURL,
			CurrentPath:    r.URL.Path,
			CurrentDate:    time.Now().Format("Monday, January 2, 2006"),
			ActiveCategory: "about",
			ImageURL:       "",
			JSONLD:         "",
			User:           user,
			Team:           team,
		}

		tmpl, err := template.ParseFiles(
			"web/templates/editorial_team.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data)
	}
}

// EditorialProfileHandler serves an individual editor's page
func EditorialProfileHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/editorial/")
		if slug == "" {
			http.Redirect(w, r, "/editorial-team", http.StatusFound)
			return
		}

		member, err := dbClient.GetEditorialTeamMember(r.Context(), slug)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Fetch recent articles by this editor using author name
		articles, _ := dbClient.GetArticlesByAuthor(r.Context(), member.Name)

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
			Member         *models.EditorialTeamMember
			Articles       []models.Article
		}{
			Title:          member.Name + " - " + cfg.SiteName,
			Description:    member.Bio,
			SiteName:       cfg.SiteName,
			SiteURL:        cfg.SiteURL,
			CurrentPath:    r.URL.Path,
			CurrentDate:    time.Now().Format("Monday, January 2, 2006"),
			ActiveCategory: "about",
			ImageURL:       member.AvatarURL,
			JSONLD:         "",
			User:           user,
			Member:         member,
			Articles:       articles,
		}

		tmpl, err := template.ParseFiles(
			"web/templates/editorial_profile.html",
			"web/templates/partials/header.html",
			"web/templates/partials/navbar.html",
			"web/templates/partials/footer.html",
		)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data)
	}
}
