package models

import (
	"time"
)

// AdminUser represents an administrator user
type AdminUser struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never return password
	CreatedAt time.Time `json:"created_at"`
}

// Media represents an uploaded file (image)
type Media struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

// ArticleStatus represents the status of an article
type ArticleStatus string

const (
	StatusDraft     ArticleStatus = "draft"
	StatusPublished ArticleStatus = "published"
	StatusArchived  ArticleStatus = "archived"
)

// Extended Article struct to include status and SEO fields if not already present
// We might need to update the main Article struct in models/article.go instead,
// but for now, let's define what we need for the admin panel.
