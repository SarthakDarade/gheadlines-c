package handlers

import (
	"gheadlines/config"
	"gheadlines/db"
	"gheadlines/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// HomeData holds data for the homepage template
type HomeData struct {
	Articles              []models.Article
	FeaturedArticle       *models.Article
	Categories            []models.Category
	CurrentCategory       string
	MetaTags              template.HTML
	CurrentYear           int
	User                  *models.User
	TrendingNews          []models.TrendingNews
	LiveUpdates           []models.LiveUpdate
	SportsArticles        []models.Article
	CountryNews           []models.Article
	BusinessArticles      []models.Article
	TechnologyArticles    []models.Article
	EntertainmentArticles []models.Article
	HealthArticles        []models.Article
}

// HomeHandler handles the homepage and category pages
func HomeHandler(dbClient *db.Client, siteURL string, siteName string, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get category filter from query
		categorySlug := r.URL.Query().Get("category")
		var categoryFilter *string
		var currentCategoryName string

		if categorySlug != "" {
			// Convert slug to category name (simple mapping for now, ideally fetch from DB)
			// In a real app, you'd look up the category by slug
			catName := strings.Title(categorySlug) // Basic capitalization
			categoryFilter = &catName
			currentCategoryName = catName
		}

		// Fetch articles
		limit := 12
		offset := 0
		// Simple pagination
		currentPage := 1
		if page := r.URL.Query().Get("page"); page != "" {
			if p, err := strconv.Atoi(page); err == nil && p > 0 {
				offset = (p - 1) * limit
				currentPage = p
			}
		}

		// Count total articles for pagination
		totalArticles, _ := dbClient.CountArticles(r.Context(), categoryFilter)
		totalPages := 0
		if totalArticles > 0 {
			totalPages = (totalArticles + limit - 1) / limit // Ceiling division
		}

		// Get access token from cookie
		var accessToken string
		if cookie, err := r.Cookie("access_token"); err == nil {
			accessToken = cookie.Value
		}

		articles, err := dbClient.GetArticles(r.Context(), limit, offset, categoryFilter, accessToken)
		if err != nil {
			// Log error but try to continue with empty list
			articles = []models.Article{}
		}

		// Fetch categories for menu
		categories, _ := dbClient.GetCategories(r.Context(), accessToken)

		// Get current user
		user, _ := GetCurrentUser(r, dbClient, cfg)

		// Fetch Trending News (only for homepage)
		var trendingNews []models.TrendingNews
		if categorySlug == "" {
			var err error
			trendingNews, err = dbClient.GetTrendingNews(r.Context(), 8, accessToken)
			if err != nil {
				log.Printf("Failed to fetch trending news: %v", err)
			} else {
				log.Printf("Fetched %d trending news items", len(trendingNews))
			}
		}

		// Fetch Live Updates (only for homepage)
		var liveUpdates []models.LiveUpdate
		if categorySlug == "" {
			liveUpdates, _ = dbClient.GetLiveUpdates(r.Context(), 10, accessToken)
		}

		// Fetch Sports Articles (only for homepage)
		var sportsArticles []models.Article
		if categorySlug == "" {
			sportsCat := "Sports"
			sportsArticles, _ = dbClient.GetArticles(r.Context(), 6, 0, &sportsCat, accessToken)
		}

		// Fetch Country News (India Focus)
		var countryNews []models.Article
		if categorySlug == "" {
			// In a real app, filtering by 'India' category or tag
			// For now, we fetch distinct articles using offset
			countryNews, _ = dbClient.GetArticles(r.Context(), 4, 12, nil, accessToken)
		}

		// Fetch Business Articles
		var businessArticles []models.Article
		if categorySlug == "" {
			cat := "Business"
			businessArticles, _ = dbClient.GetArticles(r.Context(), 4, 0, &cat, accessToken)
		}

		// Fetch Technology Articles
		var technologyArticles []models.Article
		if categorySlug == "" {
			cat := "Technology"
			technologyArticles, _ = dbClient.GetArticles(r.Context(), 4, 0, &cat, accessToken)
		}

		// Fetch Entertainment Articles
		var entertainmentArticles []models.Article
		if categorySlug == "" {
			cat := "Entertainment"
			entertainmentArticles, _ = dbClient.GetArticles(r.Context(), 4, 0, &cat, accessToken)
		}

		// Fetch Health Articles
		var healthArticles []models.Article
		if categorySlug == "" {
			cat := "Health"
			healthArticles, _ = dbClient.GetArticles(r.Context(), 4, 0, &cat, accessToken)
		}

		// Generate pages slice (simple version: all pages up to a limit)
		// TODO: Implement smart pagination (1 ... 4 5 6 ... 10)
		var pages []int
		maxPagesToShow := 10
		startPage := 1
		endPage := totalPages

		if totalPages > maxPagesToShow {
			// Show a window around current page
			startPage = currentPage - 4
			if startPage < 1 {
				startPage = 1
			}
			endPage = startPage + maxPagesToShow - 1
			if endPage > totalPages {
				endPage = totalPages
				startPage = endPage - maxPagesToShow + 1
				if startPage < 1 {
					startPage = 1
				}
			}
		}

		for i := startPage; i <= endPage; i++ {
			pages = append(pages, i)
		}

		// Prepare data
		data := struct {
			SiteName              string
			SiteURL               string
			CurrentDate           string
			Articles              []models.Article
			FeaturedArticle       *models.Article
			Categories            []models.Category
			ActiveCategory        string
			CurrentCategory       string
			JSONLD                template.HTML
			Title                 string
			Description           string
			ImageURL              string
			CurrentPath           string
			User                  *models.User
			CurrentYear           int
			TrendingNews          []models.TrendingNews
			LiveUpdates           []models.LiveUpdate
			SportsArticles        []models.Article
			CountryNews           []models.Article
			BusinessArticles      []models.Article
			TechnologyArticles    []models.Article
			EntertainmentArticles []models.Article
			HealthArticles        []models.Article
			CurrentPage           int
			TotalPages            int
			Pages                 []int
		}{
			SiteName:              siteName,
			SiteURL:               siteURL,
			CurrentDate:           time.Now().Format("Monday, January 2, 2006"),
			Articles:              articles,
			Categories:            categories,
			ActiveCategory:        categorySlug,
			CurrentCategory:       currentCategoryName,
			Title:                 currentCategoryName, // Will be empty for home
			CurrentPath:           r.URL.Path,
			User:                  user,
			CurrentYear:           time.Now().Year(),
			TrendingNews:          trendingNews,
			LiveUpdates:           liveUpdates,
			SportsArticles:        sportsArticles,
			CountryNews:           countryNews,
			BusinessArticles:      businessArticles,
			TechnologyArticles:    technologyArticles,
			EntertainmentArticles: entertainmentArticles,
			HealthArticles:        healthArticles,
			CurrentPage:           currentPage,
			TotalPages:            totalPages,
			Pages:                 pages,
		}

		// Logic for Featured Article (First one)
		if len(articles) > 0 {
			data.FeaturedArticle = &articles[0]
			// If on category page, we might want to remove the featured one from the main list
			// But for simplicity, we'll keep it or let the template handle it
		}

		// Select Template
		templateFile := "web/templates/home.html"
		if categorySlug != "" {
			templateFile = "web/templates/category.html"
		}

		// Parse Templates
		funcMap := template.FuncMap{
			"add": func(a, b int) int { return a + b },
			"sub": func(a, b int) int { return a - b },
		}

		tmpl := template.New(filepath.Base(templateFile)).Funcs(funcMap)
		tmpl, err = tmpl.ParseFiles(
			templateFile,
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
