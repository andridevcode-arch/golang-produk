package v1

import (
	"net/http"

	"go-produk/internal/config"
	"go-produk/internal/database"
	"go-produk/internal/services"
	"go-produk/pkg/token"
)

// Routes mendaftarkan semua route API v1 ke HTTP mux
func Routes(mux *http.ServeMux, authService *services.AuthService, userService *services.UserService, catService *services.CategoryService, productService *services.ProductService, templateService *services.TemplateService, templateVersionService *services.TemplateVersionService, userProjectService *services.UserProjectService, publishedPageService *services.PublishedPageService, tokenMgr *token.Manager, redisClient *database.RedisClient, cfg *config.Config) {
	// Health check (publik, tidak butuh auth)
	HealthRoutes(mux)

	// Auth (publik, tidak butuh auth)
	AuthRoutes(mux, authService, redisClient, cfg.SupabaseURL, cfg.GoogleRedirectURL)

	// User (butuh auth)
	UserRoutes(mux, userService, tokenMgr, redisClient)

	// Category
	CategoryRoutes(mux, catService, tokenMgr, redisClient)

	// Product
	ProductRoutes(mux, productService, tokenMgr, redisClient)

	// Template
	TemplateRoutes(mux, templateService, tokenMgr, redisClient)
	TemplateVersionRoutes(mux, templateVersionService, tokenMgr, redisClient)

	// User project
	UserProjectRoutes(mux, userProjectService, tokenMgr, redisClient)
	PublishedPageRoutes(mux, publishedPageService, tokenMgr, redisClient)
}
