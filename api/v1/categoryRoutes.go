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

// CategoryRoutes mendaftarkan rute untuk manajemen kategori
func CategoryRoutes(mux *http.ServeMux, catService *services.CategoryService, tokenMgr *token.Manager, redisClient *database.RedisClient, cookies cookie.Settings) {
	handler := handlers.NewCategoryHandler(catService)

	// Public routes
	mux.HandleFunc("GET /api/v1/categories", handler.GetAll)
	mux.HandleFunc("GET /api/v1/categories/{id}", handler.GetByID)

	// Protected routes (Hanya Role Admin yang diizinkan untuk Create, Update, Delete)
	authMiddleware := middlewares.AuthMiddleware(tokenMgr, redisClient, cookies)
	adminMiddleware := middlewares.RequireAdmin()

	mux.Handle("POST /api/v1/categories", authMiddleware(adminMiddleware(http.HandlerFunc(handler.Create))))
	mux.Handle("PUT /api/v1/categories/{id}", authMiddleware(adminMiddleware(http.HandlerFunc(handler.Update))))
	mux.Handle("DELETE /api/v1/categories/{id}", authMiddleware(adminMiddleware(http.HandlerFunc(handler.Delete))))
}
