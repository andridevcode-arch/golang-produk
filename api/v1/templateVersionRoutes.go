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

func TemplateVersionRoutes(mux *http.ServeMux, service *services.TemplateVersionService, tokenMgr *token.Manager, redisClient *database.RedisClient, cookies cookie.Settings) {
	h := handlers.NewTemplateVersionHandler(service)
	mux.Handle("GET /api/v1/templates/{template_id}/versions", http.HandlerFunc(h.GetAll))
	mux.Handle("GET /api/v1/templates/{template_id}/versions/{version_id}", http.HandlerFunc(h.GetByID))
	auth := middlewares.AuthMiddleware(tokenMgr, redisClient, cookies)
	admin := middlewares.RequireAdmin()
	mux.Handle("POST /api/v1/templates/{template_id}/versions", auth(admin(http.HandlerFunc(h.Create))))
	mux.Handle("PATCH /api/v1/templates/{template_id}/versions/{version_id}/active", auth(admin(http.HandlerFunc(h.Activate))))
	mux.Handle("DELETE /api/v1/templates/{template_id}/versions/{version_id}", auth(admin(http.HandlerFunc(h.Delete))))
}
