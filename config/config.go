package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	SupabaseURL           string
	SupabaseKey           string
	Port                  int
	Host                  string
	CloudflareCDN         string
	CloudflareAccountID   string
	CloudflareAPIKey      string
	CloudflareImagesToken string
	SiteURL               string
	SiteName              string
	AdminUsername         string
	AdminPasswordHash     string
	JWTSecret             string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	port := getEnvAsInt("PORT", 5000)
	host := getEnv("HOST", "localhost")

	cfg := &Config{
		SupabaseURL:           getEnv("SUPABASE_URL", ""),
		SupabaseKey:           getEnv("SUPABASE_KEY", ""),
		Port:                  port,
		Host:                  host,
		CloudflareCDN:         getEnv("CLOUDFLARE_CDN_URL", ""),
		CloudflareAccountID:   getEnv("CLOUDFLARE_ACCOUNT_ID", ""),
		CloudflareAPIKey:      getEnv("CLOUDFLARE_API_KEY", ""),
		CloudflareImagesToken: getEnv("CLOUDFLARE_IMAGES_TOKEN", ""),
		SiteURL:               getEnv("SITE_URL", "https://gh.trserver.site"),
		SiteName:              getEnv("SITE_NAME", "Global Headlines"),
		AdminUsername:         getEnv("ADMIN_USERNAME", "admin"),
		// Default password is "admin" hashed with bcrypt (cost 10)
		AdminPasswordHash: getEnv("ADMIN_PASSWORD_HASH", "$2a$10$2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2.2"),
		JWTSecret:         getEnv("JWT_SECRET", "super-secret-jwt-key-change-me"),
	}

	// Make Supabase optional - use defaults if not set
	if cfg.SupabaseURL == "" {
		cfg.SupabaseURL = "https://dcqppqeydxdnbscharjn.supabase.co"
	}
	if cfg.SupabaseKey == "" {
		cfg.SupabaseKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImRjcXBwcWV5ZHhkbmJzY2hhcmpuIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjQwNzk5NjgsImV4cCI6MjA3OTY1NTk2OH0.Mxt9n-mGcvZb4QIukohuiJA4fkiOxtKgdwPVgLhaNmk"
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
