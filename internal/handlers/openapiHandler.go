package handlers

import (
	"fmt"
	"log"
	"net/http"

	"go-produk/docs/openapi"
)

func openAPISpecHandler(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spec, err := openapi.ReadFile(filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read OpenAPI spec: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		if _, err := w.Write(spec); err != nil {
			log.Printf("failed to write OpenAPI spec: %v", err)
		}
	}
}

// OpenAPIHandler serves the correct OpenAPI spec based on query parameter
func OpenAPIHandler(baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := "openapi-user.yaml"
		if r.URL.Query().Get("type") == "admin" {
			filename = "openapi-admin.yaml"
		}
		openAPISpecHandler(filename).ServeHTTP(w, r)
	}
}

// OpenAPIUserHandler serves user OpenAPI spec only
func OpenAPIUserHandler(baseURL string) http.HandlerFunc {
	return openAPISpecHandler("openapi-user.yaml")
}

// OpenAPIAdminHandler serves admin OpenAPI spec only
func OpenAPIAdminHandler(baseURL string) http.HandlerFunc {
	return openAPISpecHandler("openapi-admin.yaml")
}
