package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gheadlines/config"
	"gheadlines/db"
	"gheadlines/models"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AuthHandler handles authentication pages
func AuthHandler(cfg *config.Config, mode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Mode":        mode, // "signin" or "signup"
			"SiteURL":     cfg.SiteURL,
			"SiteName":    cfg.SiteName,
			"SupabaseURL": cfg.SupabaseURL,
			"SupabaseKey": cfg.SupabaseKey,
		}

		tmpl, err := template.ParseFiles("web/templates/auth.html")
		if err != nil {
			http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}
}

// SignInHandler handles email/password sign in
func SignInHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		// Call Supabase auth API
		authURL := fmt.Sprintf("%s/auth/v1/token?grant_type=password", cfg.SupabaseURL)
		formData := url.Values{}
		formData.Set("email", email)
		formData.Set("password", password)

		req, err := http.NewRequest("POST", authURL, strings.NewReader(formData.Encode()))
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("apikey", cfg.SupabaseKey)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			http.Error(w, fmt.Sprintf("Authentication failed: %s", string(body)), http.StatusUnauthorized)
			return
		}

		var authResp models.AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			http.Error(w, "Failed to parse response", http.StatusInternalServerError)
			return
		}

		// Set cookie with access token
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    authResp.AccessToken,
			Path:     "/",
			MaxAge:   3600 * 24 * 7, // 7 days
			HttpOnly: true,
			Secure:   strings.HasPrefix(cfg.SiteURL, "https"),
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SignUpHandler handles email/password sign up
func SignUpHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		fullName := r.FormValue("full_name")

		// Call Supabase auth API
		authURL := fmt.Sprintf("%s/auth/v1/signup", cfg.SupabaseURL)
		formData := url.Values{}
		formData.Set("email", email)
		formData.Set("password", password)
		if fullName != "" {
			formData.Set("data", fmt.Sprintf(`{"full_name":"%s"}`, fullName))
		}

		req, err := http.NewRequest("POST", authURL, strings.NewReader(formData.Encode()))
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("apikey", cfg.SupabaseKey)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Sign up failed", http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			http.Error(w, fmt.Sprintf("Sign up failed: %s", string(body)), http.StatusBadRequest)
			return
		}

		var authResp models.AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			http.Error(w, "Failed to parse response", http.StatusInternalServerError)
			return
		}

		// Set cookie with access token
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    authResp.AccessToken,
			Path:     "/",
			MaxAge:   3600 * 24 * 7, // 7 days
			HttpOnly: true,
			Secure:   strings.HasPrefix(cfg.SiteURL, "https"),
			SameSite: http.SameSiteLaxMode,
		})

		// Create profile in public.profiles
		if authResp.User != nil {
			profile := &models.Profile{
				ID:        authResp.User.ID,
				FullName:  fullName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			go func() {
				// Use a background context for this async operation
				dbClient.CreateProfile(context.Background(), profile)
			}()
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SignOutHandler handles sign out
func SignOutHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   strings.HasPrefix(cfg.SiteURL, "https"),
			SameSite: http.SameSiteLaxMode,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// GetCurrentUser gets the current authenticated user from cookie
func GetCurrentUser(r *http.Request, dbClient *db.Client, cfg *config.Config) (*models.User, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return nil, err
	}

	// Verify token with Supabase
	userURL := fmt.Sprintf("%s/auth/v1/user", cfg.SupabaseURL)
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", cfg.SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+cookie.Value)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	// Fetch profile from public.profiles
	profile, err := dbClient.GetProfile(r.Context(), user.ID, cookie.Value)
	if err == nil && profile != nil {
		// Populate user fields from profile
		user.FullName = profile.FullName
		user.AvatarURL = profile.AvatarURL
		// Add other fields if User struct has them
	}

	return &user, nil
}

// OAuthCallbackHandler handles OAuth redirects from Supabase
func OAuthCallbackHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"SupabaseURL": cfg.SupabaseURL,
			"SupabaseKey": cfg.SupabaseKey,
		}

		tmpl, err := template.ParseFiles("web/templates/auth-callback.html")
		if err != nil {
			http.Error(w, "Failed to load callback page: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Failed to render callback page", http.StatusInternalServerError)
			return
		}
	}
}

// AuthResetPasswordHandler triggers a password reset email
func AuthResetPasswordHandler(dbClient *db.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get user from cookie to ensure they are logged in and requesting for themselves
		// (Since this is triggered from Settings page)
		user, err := GetCurrentUser(r, dbClient, cfg)
		if err != nil || user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := dbClient.ResetPasswordForEmail(r.Context(), user.Email); err != nil {
			http.Error(w, "Failed to send reset email: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reset email sent"))
	}
}
