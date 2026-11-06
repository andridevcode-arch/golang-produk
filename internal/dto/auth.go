package dto

import "go-produk/internal/models"

// RegisterRequest adalah payload untuk endpoint register
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest adalah payload untuk endpoint login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse adalah response untuk login dan register
// Token JWT tidak lagi dikirim di body, melainkan lewat HttpOnly cookie
type LoginResponse struct {
	User      models.User `json:"user"`
	ExpiresIn int         `json:"expires_in"`
}
