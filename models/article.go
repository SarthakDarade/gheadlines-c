package models

import (
	"html/template"
	"time"
)

// Article represents a news article matching Supabase schema
type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Excerpt     string    `json:"excerpt"`
	Content     string    `json:"content"`
	Category    string    `json:"category"`
	Date        string    `json:"date"`
	Author      string    `json:"author"`
	AuthorImage string    `json:"author_image"`
	AuthorBio   string    `json:"author_bio"`
	ImageURL    string    `json:"image_url"`
	ReadTime    string    `json:"read_time"`
	Views       int       `json:"views"`
	Likes       int       `json:"likes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Computed fields for UI
	Slug        string        `json:"slug,omitempty"`
	CategoryObj *Category     `json:"category_obj,omitempty"`
	ContentHTML template.HTML `json:"-"` // For template rendering
}

// Category represents a news category
type Category struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Color       string `json:"color"` // For UI styling
}

// TrendingNews represents a trending news item
type TrendingNews struct {
	ID        string    `json:"id,omitempty"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	ImageURL  string    `json:"image_url"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// LiveUpdate represents a live news update
type LiveUpdate struct {
	ID        string    `json:"id,omitempty"`
	Headline  string    `json:"headline"`
	Source    string    `json:"source"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// CareerApplication represents a job application
type CareerApplication struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ResumeURL string    `json:"resume_url"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// NewsletterSubscriber represents a newsletter subscriber
type NewsletterSubscriber struct {
	ID        string    `json:"id,omitempty"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// ContactMessage represents a contact form submission
type ContactMessage struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// BreakingNews represents a breaking news item
type BreakingNews struct {
	ID        string    `json:"id,omitempty"`
	Headline  string    `json:"headline"`
	URL       string    `json:"url"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// EditorialTeamMember represents a member of the editorial team
type EditorialTeamMember struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name"`
	Role        string         `json:"role"`
	Bio         string         `json:"bio"`
	AvatarURL   string         `json:"avatar_url"`
	SocialLinks map[string]any `json:"social_links"` // Handle JSONB
	Slug        string         `json:"slug"`
	CreatedAt   time.Time      `json:"created_at,omitempty"`
}
