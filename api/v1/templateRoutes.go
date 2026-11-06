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

func TemplateRoutes(mux *http.ServeMux, service *services.TemplateService, tokenMgr *token.Manager, redisClient *database.RedisClient, cookies cookie.Settings) {
	h := handlers.NewTemplateHandler(service)
	mux.Handle("GET /api/v1/templates", http.HandlerFunc(h.GetAll))
	mux.Handle("GET /api/v1/templates/{id}", http.HandlerFunc(h.GetByID))
	auth := middlewares.AuthMiddleware(tokenMgr, redisClient, cookies)
	admin := middlewares.RequireAdmin()
	mux.Handle("POST /api/v1/templates", auth(admin(http.HandlerFunc(h.Create))))
	mux.Handle("PUT /api/v1/templates/{id}", auth(admin(http.HandlerFunc(h.Update))))
	mux.Handle("DELETE /api/v1/templates/{id}", auth(admin(http.HandlerFunc(h.Delete))))
}
