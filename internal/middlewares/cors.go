package middlewares

import (
	"net/http"
	"strings"
)

// CORSMiddleware mengizinkan request lintas domain dengan kredensial (cookie).
// Hanya origin yang terdaftar pada allowedOrigins yang diberi header CORS.
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[strings.TrimSpace(o)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Vary wajib agar cache tidak mencampur response antar origin
			w.Header().Add("Vary", "Origin")

			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			if _, ok := allowed[origin]; !ok {
				// Origin tidak terdaftar: tidak diberi header CORS,
				// sehingga browser akan memblokir response-nya sendiri.
				next.ServeHTTP(w, r)
				return
			}

			// Origin spesifik, bukan "*", karena Allow-Credentials tidak
			// kompatibel dengan wildcard.
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
			w.Header().Set("Access-Control-Max-Age", "600")

			// Preflight dibalas lebih dulu tanpa diteruskan ke handler
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
