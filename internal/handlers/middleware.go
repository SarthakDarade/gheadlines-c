package handlers

import (
	"context"
	"gheadlines/internal/utils"
	"net/http"
	"time"
)

type contextKey string

const CSRFTokenKey contextKey = "csrf_token"

// CSRFMiddleware adds CSRF protection
func CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get or create CSRF token cookie
		cookie, err := r.Cookie("csrf_token")
		var token string

		if err != nil || cookie.Value == "" {
			// Generate new token
			token, err = utils.GenerateRandomString(32)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// Set cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    token,
				Path:     "/",
				HttpOnly: false, // JS needs to read it for fetch requests if needed, but for forms HttpOnly is safer if we inject it.
				// Actually, for Double Submit Cookie, JS often needs to read it.
				// But here we are rendering server-side templates, so we can inject it into the form.
				// So HttpOnly: true is better for security against XSS reading it, but then we must inject it into every form.
				// Let's stick to HttpOnly: true and inject via template.
				Expires:  time.Now().Add(24 * time.Hour),
				SameSite: http.SameSiteStrictMode,
			})
		} else {
			token = cookie.Value
		}

		// 2. Validate on unsafe methods
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete || r.Method == http.MethodPatch {
			// Check header
			requestToken := r.Header.Get("X-CSRF-Token")
			if requestToken == "" {
				// Check form value
				requestToken = r.FormValue("csrf_token")
			}

			if requestToken == "" || requestToken != token {
				http.Error(w, "CSRF token mismatch", http.StatusForbidden)
				return
			}
		}

		// 3. Add token to context for templates
		ctx := context.WithValue(r.Context(), CSRFTokenKey, token)
		next(w, r.WithContext(ctx))
	}
}
