package models

import "time"

// UserProfile merepresentasikan data profil tambahan user di tabel user_profiles
type UserProfile struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	DisplayName string    `json:"display_name"`
	Phone       *string   `json:"phone"`
	Bio         *string   `json:"bio"`
	AvatarURL   *string   `json:"avatar_url"`
	Location    *string   `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
