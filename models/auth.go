package models

// User represents an authenticated user
type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	AvatarURL   string `json:"avatar_url"`
	AccessToken string `json:"access_token,omitempty"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

