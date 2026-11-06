package models

import "time"

// User merepresentasikan tabel users di database
type User struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Username        *string    `json:"username,omitempty"`
	Email           string     `json:"email"`
	PasswordHash    string     `json:"-"`
	Role            string     `json:"role"`
	IsActive        bool       `json:"is_active"`
	IsDeleted       bool       `json:"is_deleted"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
