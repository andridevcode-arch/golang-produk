package v1

import (
	"net/http"

	"go-produk/internal/config"
	"go-produk/internal/database"
	"go-produk/internal/handlers"
	"go-produk/internal/services"
)

// AuthRoutes mendaftarkan route autentikasi ke mux
func AuthRoutes(mux *http.ServeMux, authService *services.AuthService, redisClient *database.RedisClient, cfg *config.Config) {
	authHandler := handlers.NewAuthHandler(authService, redisClient, cfg.Cookie, cfg.FrontendURL, cfg.SupabaseURL, cfg.GoogleRedirectURL)

	// Login
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	// Register
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)

	// Logout
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)

	// Google OAuth via Supabase Auth
	mux.HandleFunc("GET /api/v1/auth/google/login", authHandler.GoogleLogin)
	mux.HandleFunc("GET /api/v1/auth/google/callback", authHandler.GoogleCallback)
}
