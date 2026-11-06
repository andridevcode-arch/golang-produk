package v1

import (
	"net/http"

	"go-produk/internal/database"
	"go-produk/internal/handlers"
	"go-produk/internal/middlewares"
	"go-produk/internal/services"
	"go-produk/pkg/token"
)

func UserProjectRoutes(mux *http.ServeMux, service *services.UserProjectService, tokenMgr *token.Manager, redisClient *database.RedisClient) {
	h := handlers.NewUserProjectHandler(service)
	auth := middlewares.AuthMiddleware(tokenMgr, redisClient)

	mux.Handle("POST /api/v1/projects", auth(http.HandlerFunc(h.Create)))
	mux.Handle("GET /api/v1/projects", auth(http.HandlerFunc(h.List)))
	mux.Handle("GET /api/v1/projects/{id}", auth(http.HandlerFunc(h.GetByID)))
	mux.Handle("PATCH /api/v1/projects/{id}", auth(http.HandlerFunc(h.Update)))
	mux.Handle("PATCH /api/v1/projects/{id}/status", auth(http.HandlerFunc(h.UpdateStatus)))
	mux.Handle("DELETE /api/v1/projects/{id}", auth(http.HandlerFunc(h.Delete)))
}
