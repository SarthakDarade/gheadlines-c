package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"gheadlines/config"
	"gheadlines/db"
	"gheadlines/handlers"
	internalHandlers "gheadlines/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: %v. Using dummy data for MVP.", err)
		// Set defaults for MVP
		cfg = &config.Config{
			Port:     5000,
			Host:     "localhost",
			SiteURL:  "http://localhost:5000",
			SiteName: "Global Headlines",
		}
	}

	// Initialize database client
	dbClient := db.NewClient(cfg)

	// Setup routes
	mux := http.NewServeMux()

	// Homepage
	mux.HandleFunc("/", handlers.HomeHandler(dbClient, cfg.SiteURL, cfg.SiteName, cfg))

	// Article page
	mux.HandleFunc("/article/", handlers.ArticleHandler(dbClient, cfg.SiteURL, cfg.SiteName, cfg))

	// Category page (reuse home handler with category filter)
	mux.HandleFunc("/category/", func(w http.ResponseWriter, r *http.Request) {
		// Extract category slug from path
		categorySlug := strings.TrimPrefix(r.URL.Path, "/category/")
		if categorySlug == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		// Add category to query and use home handler
		q := r.URL.Query()
		q.Set("category", categorySlug)
		r.URL.RawQuery = q.Encode()
		handlers.HomeHandler(dbClient, cfg.SiteURL, cfg.SiteName, cfg)(w, r)
	})

	// Sitemap
	mux.HandleFunc("/sitemap.xml", handlers.SitemapHandler(dbClient, cfg.SiteURL))
	mux.HandleFunc("/sitemap-news.xml", handlers.NewsSitemapHandler(dbClient, cfg.SiteURL, cfg.SiteName))
	mux.HandleFunc("/news-sitemap.xml", handlers.NewsSitemapHandler(dbClient, cfg.SiteURL, cfg.SiteName)) // Alias

	// RSS Feed
	mux.HandleFunc("/rss", handlers.RSSHandler(dbClient, cfg.SiteURL, cfg.SiteName))
	mux.HandleFunc("/rss.xml", handlers.RSSHandler(dbClient, cfg.SiteURL, cfg.SiteName))

	// Robots.txt
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/robots.txt")
	})

	// Static Pages
	mux.HandleFunc("/about", internalHandlers.StaticPageHandler("about", dbClient, cfg))
	mux.HandleFunc("/contact", internalHandlers.StaticPageHandler("contact", dbClient, cfg))
	mux.HandleFunc("/privacy", internalHandlers.StaticPageHandler("privacy", dbClient, cfg))
	mux.HandleFunc("/terms", internalHandlers.StaticPageHandler("terms", dbClient, cfg))

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Real-time updates (SSE)
	mux.HandleFunc("/api/realtime", handlers.RealtimeHandler(dbClient))

	// Latest articles API (for polling fallback)
	mux.HandleFunc("/api/articles/latest", handlers.LatestArticlesHandler(dbClient))

	// User Pages
	mux.HandleFunc("/user/profile", handlers.UserProfileHandler(dbClient, cfg))
	mux.HandleFunc("/user/edit", handlers.UserEditProfileHandler(dbClient, cfg))
	mux.HandleFunc("/user/dashboard", handlers.UserDashboardHandler(dbClient, cfg))
	mux.HandleFunc("/user/settings", handlers.UserSettingsHandler(dbClient, cfg))

	// Careers
	mux.HandleFunc("/careers", handlers.CareersHandler(dbClient, cfg.SiteURL, cfg.SiteName, cfg))
	mux.HandleFunc("/api/careers/apply", handlers.SubmitCareerApplicationHandler(dbClient, cfg))

	// Forms API
	mux.HandleFunc("/api/newsletter/subscribe", handlers.NewsletterSubscribeHandler(dbClient))
	mux.HandleFunc("/api/contact", handlers.ContactFormHandler(dbClient))
	mux.HandleFunc("/api/live-updates", handlers.LiveUpdatesAPIHandler(dbClient))
	mux.HandleFunc("/api/market-data", handlers.MarketDataHandler())

	// Search
	mux.HandleFunc("/search", handlers.SearchHandler(dbClient, cfg))

	// Subscription
	mux.HandleFunc("/subscribe", handlers.SubscriptionHandler(dbClient, cfg))

	// Feature 1: Breaking News API
	mux.HandleFunc("/api/breaking-news", handlers.BreakingNewsAPIHandler(dbClient))

	// Feature 3: Editorial Team
	mux.HandleFunc("/editorial-team", handlers.EditorialTeamHandler(dbClient, cfg))
	mux.HandleFunc("/editorial/", handlers.EditorialProfileHandler(dbClient, cfg))

	// Feature 6: User Profile Update
	mux.HandleFunc("/user/update", handlers.UserUpdateProfileHandler(dbClient, cfg))

	// Authentication routes
	mux.HandleFunc("/signin", handlers.AuthHandler(cfg, "signin"))
	mux.HandleFunc("/signup", handlers.AuthHandler(cfg, "signup"))
	mux.HandleFunc("/auth/signin", handlers.SignInHandler(cfg))
	mux.HandleFunc("/auth/signup", handlers.SignUpHandler(dbClient, cfg))
	mux.HandleFunc("/auth/signout", handlers.SignOutHandler(cfg))
	mux.HandleFunc("/auth/callback", handlers.OAuthCallbackHandler(cfg))
	mux.HandleFunc("/auth/reset-password", handlers.AuthResetPasswordHandler(dbClient, cfg))

	// Admin Routes (New /adm path)
	// Login/Logout
	mux.HandleFunc("/adm/login", internalHandlers.CSRFMiddleware(internalHandlers.AdminLoginHandler(cfg)))
	mux.HandleFunc("/adm/logout", internalHandlers.AdminLogoutHandler())

	// Protected Admin Routes
	// Combine Auth and CSRF middleware
	protected := func(h http.HandlerFunc) http.HandlerFunc {
		return internalHandlers.AuthMiddleware(cfg, internalHandlers.CSRFMiddleware(h))
	}

	mux.HandleFunc("/adm/dashboard", protected(internalHandlers.AdminDashboardHandler(dbClient, cfg)))
	mux.HandleFunc("/adm/new", protected(internalHandlers.AdminNewArticleHandler(dbClient, cfg)))
	mux.HandleFunc("/adm/edit/", protected(internalHandlers.AdminEditArticleHandler(dbClient, cfg)))
	mux.HandleFunc("/adm/save", protected(internalHandlers.AdminSaveArticleHandler(dbClient, cfg)))
	mux.HandleFunc("/adm/upload", protected(internalHandlers.UploadImageHandler(cfg)))

	// Redirect old admin routes
	mux.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		newPath := strings.Replace(r.URL.Path, "/admin", "/adm", 1)
		http.Redirect(w, r, newPath, http.StatusMovedPermanently)
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware
	handler := applyMiddleware(mux)

	// Create server - bind to all interfaces for better compatibility
	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on http://localhost:%d", cfg.Port)
		log.Printf("📰 Global Headlines - World-Class SEO News Website")
		log.Printf("🌐 Access at: http://127.0.0.1:%d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// gzipResponseWriter wraps http.ResponseWriter with gzip compression
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func (g *gzipResponseWriter) Close() error {
	return g.writer.Close()
}

// applyMiddleware applies performance and security middleware
func applyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Performance headers
		if strings.HasPrefix(r.URL.Path, "/static/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else if strings.HasSuffix(r.URL.Path, ".xml") || strings.HasSuffix(r.URL.Path, ".txt") {
			w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hour for sitemap/robots
		} else {
			w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes for HTML
		}

		// Preconnect hints for CDN
		w.Header().Add("Link", "<https://cdn.cloudflare.com>; rel=preconnect")
		w.Header().Add("Link", "<https://images.unsplash.com>; rel=preconnect")
		w.Header().Add("Link", "<https://cdn.tailwindcss.com>; rel=preconnect")

		// Gzip compression for text-based content
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") &&
			(strings.Contains(r.Header.Get("Content-Type"), "text") ||
				strings.HasSuffix(r.URL.Path, ".html") ||
				strings.HasSuffix(r.URL.Path, ".xml") ||
				strings.HasSuffix(r.URL.Path, ".json")) {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")

			gz := gzip.NewWriter(w)
			defer gz.Close()

			gzw := &gzipResponseWriter{
				ResponseWriter: w,
				writer:         gz,
			}
			next.ServeHTTP(gzw, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
