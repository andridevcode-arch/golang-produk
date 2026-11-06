package middlewares

import (
	"net/http"

	"go-produk/pkg/response"
)

// RequireRole memeriksa apakah user memiliki salah satu role yang diizinkan
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRole(r)
			if userRole == "" {
				response.RespondUnauthorized(w, "identitas role user tidak ditemukan")
				return
			}

			for _, allowed := range allowedRoles {
				if userRole == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.RespondForbidden(w, "akses ditolak, hanya role yang diizinkan yang dapat melakukan aksi ini")
		})
	}
}

// RequireAdmin adalah shortcut untuk memastikan hanya role admin yang dapat mengakses
func RequireAdmin() func(http.Handler) http.Handler {
	return RequireRole("admin", "superadmin")
}
