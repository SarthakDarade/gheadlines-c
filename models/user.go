package models

import "time"

// Profile represents a user profile in public.profiles
type Profile struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	AvatarURL  string    `json:"avatar_url"`
	Website    string    `json:"website"`
	Location   string    `json:"location"`
	Bio        string    `json:"bio"`
	Occupation string    `json:"occupation"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
