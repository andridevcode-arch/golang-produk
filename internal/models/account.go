package models

import "time"

// AuthAccount merepresentasikan data autentikasi pihak ketiga maupun email di tabel auth_accounts
type AuthAccount struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	Provider          string    `json:"provider"`
	ProviderAccountID string    `json:"provider_account_id"`
	Email             string    `json:"email"`
	AvatarURL         *string   `json:"avatar_url"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
