package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gheadlines/config"
	"gheadlines/models"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client handles Supabase database operations
type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewClient creates a new Supabase client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.SupabaseURL,
		apiKey:  cfg.SupabaseKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetArticles fetches articles from Supabase matching new schema
func (c *Client) GetArticles(ctx context.Context, limit int, offset int, category *string, accessToken string) ([]models.Article, error) {
	url := fmt.Sprintf("%s/rest/v1/articles?select=*&order=created_at.desc&limit=%d&offset=%d", c.baseURL, limit, offset)

	if category != nil && *category != "" {
		url += fmt.Sprintf("&category=eq.%s", *category)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("Supabase Request Error: %v\n", err)
		return nil, err
	}

	// Retry with anon key if JWT is expired (401)
	if resp.StatusCode == http.StatusUnauthorized && accessToken != "" {
		resp.Body.Close() // Close previous body
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		resp, err = c.client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Supabase Error Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("supabase request failed with status: %d", resp.StatusCode)
	}

	var articles []models.Article
	// Read body first to debug
	body, _ := io.ReadAll(resp.Body)
	// fmt.Printf("Supabase Response: %s\n", string(body)) // Uncomment to see full JSON

	if err := json.Unmarshal(body, &articles); err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return nil, err
	}

	// Generate slugs and category objects for articles
	for i := range articles {
		if articles[i].Slug == "" {
			articles[i].Slug = generateSlug(articles[i].Title, articles[i].ID)
		}
		if articles[i].CategoryObj == nil {
			articles[i].CategoryObj = getCategoryByName(articles[i].Category)
		}
	}

	return articles, nil
}

// CountArticles counts articles matching the category
func (c *Client) CountArticles(ctx context.Context, category *string) (int, error) {
	if c.baseURL == "" || c.apiKey == "" {
		return 0, nil
	}

	url := fmt.Sprintf("%s/rest/v1/articles", c.baseURL)

	// Prepare query params
	queryParams := []string{"select=id", "limit=0"}
	if category != nil && *category != "" {
		queryParams = append(queryParams, fmt.Sprintf("category=eq.%s", *category))
	}

	url += "?" + strings.Join(queryParams, "&")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Prefer", "count=exact")

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return 0, fmt.Errorf("failed to count articles: status %d", resp.StatusCode)
	}

	// Parse Content-Range header
	// Format: */Count or Start-End/Count
	rangeHeader := resp.Header.Get("Content-Range")
	if rangeHeader == "" {
		return 0, nil // Should ideally be an error if we expect it
	}

	parts := strings.Split(rangeHeader, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid Content-Range header: %s", rangeHeader)
	}

	var count int
	fmt.Sscanf(parts[1], "%d", &count)
	return count, nil
}

// GetArticleByID fetches a single article by ID
func (c *Client) GetArticleByID(ctx context.Context, id string, accessToken string) (*models.Article, error) {
	url := fmt.Sprintf("%s/rest/v1/articles?id=eq.%s&limit=1", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Retry with anon key if JWT is expired (401)
	if resp.StatusCode == http.StatusUnauthorized && accessToken != "" {
		resp.Body.Close()
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		resp, err = c.client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch article: status %d", resp.StatusCode)
	}

	var articles []models.Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, err
	}

	if len(articles) == 0 {
		return nil, fmt.Errorf("article not found")
	}

	article := &articles[0]
	article.Slug = generateSlug(article.Title, article.ID)
	article.CategoryObj = getCategoryByName(article.Category)
	return article, nil
}

// GetArticleBySlug fetches article by slug (searches by ID prefix in slug)
func (c *Client) GetArticleBySlug(ctx context.Context, slug string, accessToken string) (*models.Article, error) {
	// Optimization: Extract ID prefix from slug and search by ID
	// Slug format: title-slug-idprefix (8 chars)
	if len(slug) < 9 {
		// Fallback for short slugs (legacy?) - try title match on recent articles
		// Or just return error
		return nil, fmt.Errorf("invalid slug format: %s", slug)
	}

	idPrefix := slug[len(slug)-8:]

	// Search for article with ID starting with this prefix
	// We fetch the full article directly as we expect only 1 match
	url := fmt.Sprintf("%s/rest/v1/articles?id=like.%s*&limit=5", c.baseURL, idPrefix)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Retry with anon key if JWT is expired (401)
	if resp.StatusCode == http.StatusUnauthorized && accessToken != "" {
		resp.Body.Close()
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		resp, err = c.client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch article: status %d", resp.StatusCode)
	}

	var articles []models.Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, err
	}

	// Find the exact match
	for i := range articles {
		// Generate slug to verify
		if articles[i].Slug == "" {
			articles[i].Slug = generateSlug(articles[i].Title, articles[i].ID)
		}
		if articles[i].CategoryObj == nil {
			articles[i].CategoryObj = getCategoryByName(articles[i].Category)
		}

		if articles[i].Slug == slug {
			return &articles[i], nil
		}
	}

	return nil, fmt.Errorf("article not found with slug: %s", slug)
}

// GetCategories fetches all categories
func (c *Client) GetCategories(ctx context.Context, accessToken string) ([]models.Category, error) {
	// If Supabase is not configured, return error
	if c.baseURL == "" || c.apiKey == "" {
		return nil, fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/categories?select=*&order=name.asc", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Retry with anon key if JWT is expired (401)
	if resp.StatusCode == http.StatusUnauthorized && accessToken != "" {
		resp.Body.Close()
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		resp, err = c.client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch categories: status %d", resp.StatusCode)
	}

	var categories []models.Category
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// SearchArticles searches for articles by title
func (c *Client) SearchArticles(ctx context.Context, query string) ([]models.Article, error) {
	if c.baseURL == "" || c.apiKey == "" {
		return nil, fmt.Errorf("supabase not configured")
	}

	// Simple search by title (ILIKE)
	// Note: URL encoding for query is important, but basic string should work for MVP
	// Better to use url.QueryEscape
	url := fmt.Sprintf("%s/rest/v1/articles?title=ilike.*%s*&order=created_at.desc&limit=20", c.baseURL, query)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed: status %d", resp.StatusCode)
	}

	var articles []models.Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, err
	}

	// Process articles (slugs, categories)
	for i := range articles {
		if articles[i].Slug == "" {
			articles[i].Slug = generateSlug(articles[i].Title, articles[i].ID)
		}
		if articles[i].CategoryObj == nil {
			articles[i].CategoryObj = getCategoryByName(articles[i].Category)
		}
	}

	return articles, nil
}

// CreateArticle creates a new article
func (c *Client) CreateArticle(ctx context.Context, article *models.Article) error {
	// If Supabase is not configured, just return nil (mock success)
	if c.baseURL == "" || c.apiKey == "" {
		return nil
	}

	url := fmt.Sprintf("%s/rest/v1/articles", c.baseURL)

	body, err := json.Marshal(article)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to create article: status %d", resp.StatusCode)
	}

	return nil
}

// UpdateArticle updates an existing article
func (c *Client) UpdateArticle(ctx context.Context, article *models.Article) error {
	if c.baseURL == "" || c.apiKey == "" {
		return nil
	}

	url := fmt.Sprintf("%s/rest/v1/articles?id=eq.%s", c.baseURL, article.ID)

	body, err := json.Marshal(article)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to update article: status %d", resp.StatusCode)
	}

	return nil
}

// DeleteArticle deletes an article
func (c *Client) DeleteArticle(ctx context.Context, id string) error {
	if c.baseURL == "" || c.apiKey == "" {
		return nil
	}

	url := fmt.Sprintf("%s/rest/v1/articles?id=eq.%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete article: status %d", resp.StatusCode)
	}

	return nil
}

// CreateProfile creates a new user profile
func (c *Client) CreateProfile(ctx context.Context, profile *models.Profile) error {
	if c.baseURL == "" || c.apiKey == "" {
		return fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/profiles", c.baseURL)

	body, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusConflict {
			return nil
		}
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create profile: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetProfile fetches a user profile by ID
func (c *Client) GetProfile(ctx context.Context, id string, accessToken string) (*models.Profile, error) {
	url := fmt.Sprintf("%s/rest/v1/profiles?id=eq.%s&limit=1", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch profile: status %d", resp.StatusCode)
	}

	var profiles []models.Profile
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		return nil, err
	}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("profile not found")
	}

	return &profiles[0], nil
}

// GetTrendingNews fetches trending news items
func (c *Client) GetTrendingNews(ctx context.Context, limit int, accessToken string) ([]models.TrendingNews, error) {
	if c.baseURL == "" || c.apiKey == "" {
		return []models.TrendingNews{}, nil
	}

	url := fmt.Sprintf("%s/rest/v1/trending_news?select=*&order=created_at.desc&limit=%d", c.baseURL, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Always use Anon Key for public data to avoid 401s from expired user tokens
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch trending news: status %d", resp.StatusCode)
	}

	var trending []models.TrendingNews
	if err := json.NewDecoder(resp.Body).Decode(&trending); err != nil {
		return nil, err
	}

	return trending, nil
}

// GetLiveUpdates fetches live news updates
func (c *Client) GetLiveUpdates(ctx context.Context, limit int, accessToken string) ([]models.LiveUpdate, error) {
	if c.baseURL == "" || c.apiKey == "" {
		return []models.LiveUpdate{}, nil
	}

	url := fmt.Sprintf("%s/rest/v1/live_updates?select=*&order=created_at.desc&limit=%d", c.baseURL, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Always use Anon Key for public updates
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch live updates: status %d", resp.StatusCode)
	}

	var updates []models.LiveUpdate
	if err := json.NewDecoder(resp.Body).Decode(&updates); err != nil {
		return nil, err
	}

	return updates, nil
}

// CreateCareerApplication submits a new career application
func (c *Client) CreateCareerApplication(ctx context.Context, app *models.CareerApplication) error {
	if c.baseURL == "" || c.apiKey == "" {
		return fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/careers_applications", c.baseURL)

	body, err := json.Marshal(app)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to submit application: status %d", resp.StatusCode)
	}

	return nil
}

// CreateNewsletterSubscriber adds a new newsletter subscriber
func (c *Client) CreateNewsletterSubscriber(ctx context.Context, email string) error {
	if c.baseURL == "" || c.apiKey == "" {
		return fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/newsletter_subscribers", c.baseURL)

	payload := map[string]string{"email": email}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 409 means conflict (duplicate email), which is fine for us
	if resp.StatusCode == http.StatusConflict {
		return nil
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to subscribe: status %d", resp.StatusCode)
	}

	return nil
}

// CreateContactMessage submits a contact form message
func (c *Client) CreateContactMessage(ctx context.Context, msg *models.ContactMessage) error {
	if c.baseURL == "" || c.apiKey == "" {
		return fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/contact_messages", c.baseURL)

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send message: status %d", resp.StatusCode)
	}

	return nil
}

// GetBreakingNews fetches distinct breaking news items
func (c *Client) GetBreakingNews(ctx context.Context, limit int) ([]models.BreakingNews, error) {
	if c.baseURL == "" || c.apiKey == "" {
		return []models.BreakingNews{}, nil
	}

	url := fmt.Sprintf("%s/rest/v1/breaking_news?select=*&order=created_at.desc&limit=%d", c.baseURL, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch breaking news: status %d", resp.StatusCode)
	}

	var news []models.BreakingNews
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return nil, err
	}

	return news, nil
}

// GetEditorialTeam fetches all editorial team members
func (c *Client) GetEditorialTeam(ctx context.Context) ([]models.EditorialTeamMember, error) {
	if c.baseURL == "" || c.apiKey == "" {
		return []models.EditorialTeamMember{}, nil
	}

	url := fmt.Sprintf("%s/rest/v1/editorial_team?select=*&order=name.asc", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch editorial team: status %d", resp.StatusCode)
	}

	var team []models.EditorialTeamMember
	if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
		return nil, err
	}

	return team, nil
}

// CreateBreakingNews creates a new breaking news item
func (c *Client) CreateBreakingNews(ctx context.Context, news *models.BreakingNews) error {
	if c.baseURL == "" || c.apiKey == "" {
		return fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/breaking_news", c.baseURL)

	body, err := json.Marshal(news)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to create breaking news: status %d", resp.StatusCode)
	}

	return nil
}

// CreateEditorialTeamMember creates a new editorial team member
func (c *Client) CreateEditorialTeamMember(ctx context.Context, member *models.EditorialTeamMember) error {
	if c.baseURL == "" || c.apiKey == "" {
		return fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/editorial_team", c.baseURL)

	body, err := json.Marshal(member)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=minimal")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to create team member: status %d", resp.StatusCode)
	}

	return nil
}

// GetEditorialTeamMember fetches a single team member by slug
func (c *Client) GetEditorialTeamMember(ctx context.Context, slug string) (*models.EditorialTeamMember, error) {
	if c.baseURL == "" || c.apiKey == "" {
		return nil, fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/editorial_team?slug=eq.%s&limit=1", c.baseURL, slug)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch team member: status %d", resp.StatusCode)
	}

	var members []models.EditorialTeamMember
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("member not found")
	}

	return &members[0], nil
}

// ResetPasswordForEmail sends a password reset email
func (c *Client) ResetPasswordForEmail(ctx context.Context, email string) error {
	if c.baseURL == "" || c.apiKey == "" {
		return fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/auth/v1/recover", c.baseURL)

	payload := map[string]string{"email": email}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send reset email: status %d", resp.StatusCode)
	}

	return nil
}

// UpdateUserProfile updates a user's profile
func (c *Client) UpdateUserProfile(ctx context.Context, profile *models.Profile, accessToken string) error {
	if c.baseURL == "" || c.apiKey == "" {
		return fmt.Errorf("supabase not configured")
	}

	url := fmt.Sprintf("%s/rest/v1/profiles?id=eq.%s", c.baseURL, profile.ID)

	// Use a map to avoid sending ID and CreatedAt/UpdatedAt which might be read-only or cause RLS issues
	updates := map[string]interface{}{
		"full_name":  profile.FullName,
		"username":   profile.Username,
		"bio":        profile.Bio,
		"occupation": profile.Occupation,
		"location":   profile.Location,
		"website":    profile.Website,
		"avatar_url": profile.AvatarURL,
		"updated_at": time.Now().Format(time.RFC3339),
	}

	body, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("apikey", c.apiKey)
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to update profile: status %d", resp.StatusCode)
	}

	return nil
}
