package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "go-produk/api/v1"
	"go-produk/internal/database"
	"go-produk/internal/handlers"
	"go-produk/internal/middlewares"
)

// SetupRouter membuat HTTP handler dengan route yang terdaftar
func (a *App) SetupRouter() http.Handler {
	mux := http.NewServeMux()

	// Route: / (halaman utama)
	mux.HandleFunc("/", a.handleHome)

	// Route: /health (health check semua database - publik)
	mux.HandleFunc("/health", a.handleHealth)

	// Route: /openapi.yaml (OpenAPI spec dinamis - deprecated)
	mux.Handle("GET /openapi.yaml", handlers.OpenAPIHandler(""))

	// Route: /openapi/user.yaml (User API spec)
	mux.Handle("GET /openapi/user.yaml", handlers.OpenAPIUserHandler(""))

	// Route: /openapi/admin.yaml (Admin API spec)
	mux.Handle("GET /openapi/admin.yaml", handlers.OpenAPIAdminHandler(""))

	// Route: /docs (interactive documentation dengan Scalar-Go)
	// /docs/user  - User API documentation
	mux.Handle("GET /docs/user", handlers.DocsHandler("", "user"))
	// /docs/admin - Admin API documentation
	mux.Handle("GET /docs/admin", handlers.DocsHandler("", "admin"))

	// Register semua route API v1
	v1.Routes(mux, a.AuthService, a.UserService, a.CategoryService, a.ProductService, a.TemplateService, a.TemplateVersionService, a.UserProjectService, a.PublishedPageService, a.TokenMgr, a.Redis, a.Config)
	if a.PublishedPageService != nil {
		v1.PublicPublishedPageRoutes(mux, a.PublishedPageService)
	}

	// Middleware: recovery (panic recovery)
	handler := recoveryMiddleware(mux)
	// Middleware: proteksi CSRF (validasi Origin/Referer untuk method mutasi)
	handler = middlewares.CSRFMiddleware(a.Config.AllowedOrigins)(handler)
	// Middleware: CORS ber-kredensial (berada di luar CSRF agar preflight OPTIONS dibalas lebih dulu)
	handler = middlewares.CORSMiddleware(a.Config.AllowedOrigins)(handler)
	// Middleware: request ID
	handler = requestIDMiddleware(handler)
	// Middleware: logger
	handler = loggerMiddleware(handler)

	return handler
}

// handleHome meresep request ke root path
func (a *App) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Welcome to Go Produk API",
		"status":    "running",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// handleHealth memeriksa status koneksi semua database
func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result := map[string]interface{}{
		"status":    "healthy",
		"databases": map[string]interface{}{},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	// Postgres (optional)
	if a.Supabase != nil {
		result["databases"].(map[string]interface{})["postgres"] = checkDBHealth(ctx, a.Supabase.Health)
	} else {
		result["databases"].(map[string]interface{})["postgres"] = map[string]string{
			"name":   "postgres",
			"status": "down",
			"error":  "postgres tidak dikonfigurasi",
		}
		result["status"] = "degraded"
	}

	// Redis (optional)
	if a.Redis != nil {
		result["databases"].(map[string]interface{})["redis"] = checkDBHealth(ctx, a.Redis.Health)
	} else {
		result["databases"].(map[string]interface{})["redis"] = map[string]string{
			"name":   "redis",
			"status": "down",
			"error":  "redis tidak dikonfigurasi",
		}
		result["status"] = "degraded"
	}

	// MongoDB (optional)
	if a.Mongo != nil {
		result["databases"].(map[string]interface{})["mongo"] = checkDBHealth(ctx, a.Mongo.Health)
	} else {
		result["databases"].(map[string]interface{})["mongo"] = map[string]string{
			"name":   "mongo",
			"status": "down",
			"error":  "mongo tidak dikonfigurasi",
		}
		result["status"] = "degraded"
	}

	w.Header().Set("Content-Type", "application/json")
	if result["status"] == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(result)
}

// checkDBHealth menjalankan health check dan mengembalikan HealthStatus
func checkDBHealth(ctx context.Context, healthFn func(context.Context) database.HealthStatus) interface{} {
	return healthFn(ctx)
}

// recoveryMiddleware menangani panic dan mengembalikan response 500
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v", rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// requestIDMiddleware menambahkan header X-Request-ID
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())
		r = r.WithContext(context.WithValue(r.Context(), "requestID", requestID))
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

// loggerMiddleware mencatat request masuk
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// StartServer menjalankan HTTP server dengan graceful shutdown
func (a *App) StartServer(addr string) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      a.SetupRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		a.Shutdown(shutdownCtx)
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Printf("server running on %s", addr)
	return server.ListenAndServe()
}
