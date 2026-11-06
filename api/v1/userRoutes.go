package v1

import (
	"net/http"

	"go-produk/internal/database"
	"go-produk/internal/handlers"
	"go-produk/internal/middlewares"
	"go-produk/internal/services"
	"go-produk/pkg/cookie"
	"go-produk/pkg/token"
)

// UserRoutes mendaftarkan route user ke mux
func UserRoutes(mux *http.ServeMux, userSvc *services.UserService, tokenMgr *token.Manager, redisClient *database.RedisClient, cookies cookie.Settings) {
	authMW := middlewares.AuthMiddleware(tokenMgr, redisClient, cookies)
	userHandler := handlers.NewUserHandler(userSvc)

	// Get current user profile
	mux.Handle("GET /api/v1/users/me", authMW(http.HandlerFunc(userHandler.GetProfile)))

	// Update user profile
	mux.Handle("PUT /api/v1/users/me", authMW(http.HandlerFunc(userHandler.UpdateProfile)))
}
