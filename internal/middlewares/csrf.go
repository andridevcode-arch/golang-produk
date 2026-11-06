package middlewares

import (
	"net/http"
	"net/url"
	"strings"

	"go-produk/pkg/response"
)

// CSRFMiddleware menolak request mutasi yang berasal dari origin di luar daftar.
// Middleware ini diperlukan karena cookie SameSite=None tetap terkirim
// pada request lintas situs.
func CSRFMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[strings.TrimSuffix(strings.TrimSpace(o), "/")] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Method aman langsung diteruskan
			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
				next.ServeHTTP(w, r)
				return
			}

			origin := strings.TrimSuffix(r.Header.Get("Origin"), "/")
			if origin == "" {
				origin = refererOrigin(r.Header.Get("Referer"))
			}

			if origin == "" {
				response.RespondForbidden(w, "origin request tidak dapat diverifikasi")
				return
			}

			if _, ok := allowed[origin]; !ok {
				response.RespondForbidden(w, "origin request tidak diizinkan")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// refererOrigin mengambil bagian scheme://host dari header Referer
func refererOrigin(referer string) string {
	if referer == "" {
		return ""
	}
	u, err := url.Parse(referer)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	return u.Scheme + "://" + u.Host
}
