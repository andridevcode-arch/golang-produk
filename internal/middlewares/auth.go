package middlewares

import (
	"context"
	"net/http"
	"strings"

	"go-produk/internal/database"
	"go-produk/pkg/response"
	"go-produk/pkg/token"
)

// contextKey adalah tipe kunci untuk context
type contextKey string

const (
	userIDKey   contextKey = "user_id"
	userRoleKey contextKey = "user_role"
	tokenKey    contextKey = "raw_token"
)

// GetUserID mengambil userID dari context request
func GetUserID(r *http.Request) string {
	val := r.Context().Value(userIDKey)
	if val == nil {
		return ""
	}
	return val.(string)
}

// GetUserRole mengambil userRole dari context request
func GetUserRole(r *http.Request) string {
	val := r.Context().Value(userRoleKey)
	if val == nil {
		return ""
	}
	return val.(string)
}

// GetRawToken mengambil raw token string dari context request
func GetRawToken(r *http.Request) string {
	val := r.Context().Value(tokenKey)
	if val == nil {
		return ""
	}
	return val.(string)
}

// AuthMiddleware adalah middleware untuk validasi JWT token dan pengecekan blacklist di Redis
func AuthMiddleware(tokenMgr *token.Manager, redisClient *database.RedisClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil token dari header Authorization: Bearer <token>
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.RespondUnauthorized(w, "authorization header kosong")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.RespondUnauthorized(w, "format authorization header salah, gunakan: Bearer <token>")
				return
			}

			tokenStr := strings.TrimSpace(parts[1])
			if tokenStr == "" {
				response.RespondUnauthorized(w, "token kosong")
				return
			}

			// Cek apakah token ada di blacklist (jika Redis aktif)
			if redisClient != nil {
				blacklisted, err := redisClient.IsBlacklisted(r.Context(), tokenStr)
				if err == nil && blacklisted {
					response.RespondUnauthorized(w, "token sudah tidak valid atau telah logout")
					return
				}
			}

			// Validasi token
			claims, err := tokenMgr.Validate(tokenStr)
			if err != nil {
				response.RespondUnauthorized(w, "token tidak valid atau sudah kadaluarsa")
				return
			}

			// Simpan userID, role, dan token ke context
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, userRoleKey, claims.Role)
			ctx = context.WithValue(ctx, tokenKey, tokenStr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
