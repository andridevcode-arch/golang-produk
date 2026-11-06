package dto

import (
	"time"

	"go-produk/internal/models"
)

// UserDetailResponse merepresentasikan seluruh agregasi data user, profil, dan daftar akunnya
type UserDetailResponse struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Username        *string              `json:"username,omitempty"`
	Email           string               `json:"email"`
	Role            string               `json:"role"`
	IsActive        bool                 `json:"is_active"`
	EmailVerifiedAt *time.Time           `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Profile         *models.UserProfile  `json:"profile"`
	Accounts        []models.AuthAccount `json:"accounts"`
}

// UpdateUserRequest adalah payload untuk memperbarui data nama user & user profile (email dan accounts bersifat read-only)
type UpdateUserRequest struct {
	Name        string  `json:"name"`
	Username    *string `json:"username"`
	DisplayName string  `json:"display_name"`
	Phone       *string `json:"phone"`
	Bio         *string `json:"bio"`
	AvatarURL   *string `json:"avatar_url"`
	Location    *string `json:"location"`
}
