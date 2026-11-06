package handlers

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAPISpecHandlers(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		target      string
		wantContent string
	}{
		{
			name:        "user",
			handler:     OpenAPIUserHandler(""),
			target:      "/openapi/user.yaml",
			wantContent: "/api/v1/auth/register:",
		},
		{
			name:        "admin",
			handler:     OpenAPIAdminHandler(""),
			target:      "/openapi/admin.yaml",
			wantContent: "/api/v1/categories:",
		},
		{
			name:        "legacy user default",
			handler:     OpenAPIHandler(""),
			target:      "/openapi.yaml",
			wantContent: "/api/v1/auth/register:",
		},
		{
			name:        "legacy admin query",
			handler:     OpenAPIHandler(""),
			target:      "/openapi.yaml?type=admin",
			wantContent: "/api/v1/categories:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			tt.handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tt.target, nil))

			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
			}
			if got := recorder.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/yaml") {
				t.Fatalf("Content-Type = %q, want application/yaml", got)
			}
			if !strings.Contains(recorder.Body.String(), tt.wantContent) {
				t.Fatalf("response does not contain %q", tt.wantContent)
			}
			if !strings.Contains(recorder.Body.String(), "  - url: /") {
				t.Fatal("OpenAPI spec server URL should be relative to the current host")
			}
		})
	}
}

func TestDocsHandlerRendersScalarWithSpecURL(t *testing.T) {
	recorder := httptest.NewRecorder()
	DocsHandler("", "user").ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "http://localhost:8025/docs/user", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", got)
	}
	for _, want := range []string{
		"Go Produk API - User",
		"http://localhost:8025/openapi/user.yaml",
		`"servers":[{"url":"http://localhost:8025","description":"Current server"}`,
	} {
		if !strings.Contains(recorder.Body.String(), want) {
			t.Fatalf("Scalar page does not contain %q", want)
		}
	}
}

func TestDocsHandlerUsesRequestDomain(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "https://api.example.com/docs/user", nil)
	request.TLS = &tls.ConnectionState{}
	DocsHandler("", "user").ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	for _, want := range []string{
		"https://api.example.com/openapi/user.yaml",
		`"servers":[{"url":"https://api.example.com","description":"Current server"}`,
	} {
		if !strings.Contains(recorder.Body.String(), want) {
			t.Fatalf("Scalar page does not contain request domain %q", want)
		}
	}
}
