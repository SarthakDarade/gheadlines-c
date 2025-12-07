package handlers

import (
	"gheadlines/config"
	"gheadlines/internal/utils"
	"html/template"
	"net/http"
	"time"
)

// AdminLoginHandler handles the admin login page and submission
func AdminLoginHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Check if already logged in
			cookie, err := r.Cookie("admin_token")
			if err == nil {
				if _, err := utils.ValidateJWT(cookie.Value, cfg.JWTSecret); err == nil {
					http.Redirect(w, r, "/adm/dashboard", http.StatusSeeOther)
					return
				}
			}

			// Render login page
			tmpl, err := template.ParseFiles("web/templates/admin/login.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data := struct {
				CSRFToken string
				Error     string
			}{
				CSRFToken: r.Context().Value(CSRFTokenKey).(string),
				Error:     r.URL.Query().Get("error"),
			}

			tmpl.Execute(w, data)
			return
		}

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			// Simple check against config credentials
			// In a real app, you'd check against DB
			if username != cfg.AdminUsername || !utils.CheckPasswordHash(password, cfg.AdminPasswordHash) {
				// For MVP, if hash check fails, check plain text (for initial setup convenience)
				// WARNING: Remove this in production
				if username != cfg.AdminUsername || password != "admin" {
					http.Redirect(w, r, "/adm/login?error=invalid_credentials", http.StatusSeeOther)
					return
				}
			}

			// Generate JWT
			token, err := utils.GenerateJWT(username, cfg.JWTSecret, time.Now().Add(24*time.Hour))
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Set cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "admin_token",
				Value:    token,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				Path:     "/",
				Secure:   r.TLS != nil, // Set Secure if using HTTPS
				SameSite: http.SameSiteStrictMode,
			})

			http.Redirect(w, r, "/adm/dashboard", http.StatusSeeOther)
		}
	}
}

// AdminLogoutHandler clears the session
func AdminLogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})
		http.Redirect(w, r, "/adm/login", http.StatusSeeOther)
	}
}

// AuthMiddleware protects admin routes
func AuthMiddleware(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("admin_token")
		if err != nil {
			http.Redirect(w, r, "/adm/login", http.StatusSeeOther)
			return
		}

		claims, err := utils.ValidateJWT(cookie.Value, cfg.JWTSecret)
		if err != nil {
			http.Redirect(w, r, "/adm/login", http.StatusSeeOther)
			return
		}

		// You could add claims to context here if needed
		_ = claims

		next(w, r)
	}
}
