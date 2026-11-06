package v1

import (
	"encoding/json"
	"net/http"
)

// HealthRoutes mendaftarkan route health check khusus API ke mux
func HealthRoutes(mux *http.ServeMux) {
	// Route untuk API health
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "go-produk-api",
			"version": "v1",
		})
	})
}
