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
	AuthRoutes(mux, authService, redisClient, cfg)

	// User (butuh auth)
	UserRoutes(mux, userService, tokenMgr, redisClient, cfg.Cookie)

	// Category
	CategoryRoutes(mux, catService, tokenMgr, redisClient, cfg.Cookie)

	// Product
	ProductRoutes(mux, productService, tokenMgr, redisClient, cfg.Cookie)

	// Template
	TemplateRoutes(mux, templateService, tokenMgr, redisClient, cfg.Cookie)
	TemplateVersionRoutes(mux, templateVersionService, tokenMgr, redisClient, cfg.Cookie)

	// User project
	UserProjectRoutes(mux, userProjectService, tokenMgr, redisClient, cfg.Cookie)
	PublishedPageRoutes(mux, publishedPageService, tokenMgr, redisClient, cfg.Cookie)
}
