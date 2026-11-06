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

// ProductRoutes mendaftarkan route products ke mux
func ProductRoutes(mux *http.ServeMux, productService *services.ProductService, tokenMgr *token.Manager, redisClient *database.RedisClient, cookies cookie.Settings) {
	handler := handlers.NewProductHandler(productService)

	authMiddleware := middlewares.AuthMiddleware(tokenMgr, redisClient, cookies)
	adminMiddleware := middlewares.RequireAdmin()

	// List products (publik)
	mux.Handle("GET /api/v1/products", http.HandlerFunc(handler.GetAll))

	// Get product by ID (publik)
	mux.Handle("GET /api/v1/products/{id}", http.HandlerFunc(handler.GetByID))

	// Create product (admin only)
	mux.Handle("POST /api/v1/products", authMiddleware(adminMiddleware(http.HandlerFunc(handler.Create))))

	// Update product (admin only)
	mux.Handle("PUT /api/v1/products/{id}", authMiddleware(adminMiddleware(http.HandlerFunc(handler.Update))))

	// Delete product (admin only)
	mux.Handle("DELETE /api/v1/products/{id}", authMiddleware(adminMiddleware(http.HandlerFunc(handler.Delete))))
}
