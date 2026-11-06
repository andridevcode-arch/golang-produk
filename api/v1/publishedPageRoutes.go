package v1

import (
	"net/http"

	"go-produk/internal/database"
	"go-produk/internal/handlers"
	"go-produk/internal/middlewares"
	"go-produk/internal/services"
	"go-produk/pkg/token"
)

func PublishedPageRoutes(mux *http.ServeMux, service *services.PublishedPageService, tokenMgr *token.Manager, redisClient *database.RedisClient) {
	h := handlers.NewPublishedPageHandler(service)
	auth := middlewares.AuthMiddleware(tokenMgr, redisClient)
	mux.Handle("PUT /api/v1/projects/{id}/page", auth(http.HandlerFunc(h.Upsert)))
	mux.Handle("GET /api/v1/projects/{id}/page", auth(http.HandlerFunc(h.GetOwner)))
	mux.Handle("PATCH /api/v1/projects/{id}/page", auth(http.HandlerFunc(h.Update)))
	mux.Handle("POST /api/v1/projects/{id}/page/publish", auth(http.HandlerFunc(h.Publish)))
	mux.Handle("POST /api/v1/projects/{id}/page/unpublish", auth(http.HandlerFunc(h.Unpublish)))
}

func PublicPublishedPageRoutes(mux *http.ServeMux, service *services.PublishedPageService) {
	h := handlers.NewPublishedPageHandler(service)
	mux.Handle("GET /{username}/{slug}", http.HandlerFunc(h.Public))
	mux.Handle("GET /{username}/{slug}/preview", http.HandlerFunc(h.Preview))
}
