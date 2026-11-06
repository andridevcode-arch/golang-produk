package handlers

import (
	"fmt"
	"net/http"
	"strings"

	scalargo "github.com/bdpiprava/scalar-go"
)

// DocsHandler renders interactive API documentation using Scalar-Go.
// It dynamically loads OpenAPI spec with current server URL.
// Parameter specType: "user" or "admin" to load respective spec.
func DocsHandler(baseURL string, specType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Determine server URL
		serverURL := baseURL
		if serverURL == "" {
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			serverURL = fmt.Sprintf("%s://%s", scheme, r.Host)
		}
		serverURL = strings.TrimRight(serverURL, "/")

		// Determine spec URL based on type
		specURL := serverURL + "/openapi/" + specType + ".yaml"

		// Generate HTML using Scalar-Go V2 with dynamic spec URL
		html, err := scalargo.NewV2(
			scalargo.WithSpecURL(specURL),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Go Produk API - "+ucfirst(specType)),
				scalargo.WithKeyValue("description", "Interactive API documentation for "+ucfirst(specType)),
			),
			scalargo.WithServers(
				scalargo.ServerOverride{URL: serverURL, Description: "Current server"},
				scalargo.ServerOverride{URL: "https://api.go-produk.com", Description: "Production"},
			),
			scalargo.WithTheme(scalargo.ThemeAlternate),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to generate docs: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	}
}

// ucfirst kapitalisasi pertama karakter string
func ucfirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
