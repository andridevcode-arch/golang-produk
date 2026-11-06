package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-produk/internal/database"
	"go-produk/internal/dto"
	"go-produk/internal/services"
	"go-produk/pkg/cookie"
	"go-produk/pkg/response"
)

// AuthHandler menangani request HTTP untuk autentikasi
type AuthHandler struct {
	authService *services.AuthService
	redisClient *database.RedisClient
	cookies     cookie.Settings
	frontendURL string
	supabaseURL string
	redirectURL string
}

// NewAuthHandler membuat instance AuthHandler baru
func NewAuthHandler(authService *services.AuthService, redisClient *database.RedisClient, cookies cookie.Settings, frontendURL, supabaseURL, redirectURL string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		redisClient: redisClient,
		cookies:     cookies,
		frontendURL: frontendURL,
		supabaseURL: supabaseURL,
		redirectURL: redirectURL,
	}
}

// Register menangani POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request tidak valid")
		return
	}

	resp, token, err := h.authService.Register(r.Context(), req)
	if err != nil {
		response.RespondBadRequest(w, err.Error())
		return
	}

	// Token dikirim lewat HttpOnly cookie, bukan lewat body response
	h.cookies.Set(w, token, resp.ExpiresIn)

	response.RespondCreated(w, resp)
}

// Login menangani POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondBadRequest(w, "format request tidak valid")
		return
	}

	resp, token, err := h.authService.Login(r.Context(), req)
	if err != nil {
		response.RespondUnauthorized(w, err.Error())
		return
	}

	// Token dikirim lewat HttpOnly cookie, bukan lewat body response
	h.cookies.Set(w, token, resp.ExpiresIn)

	response.RespondSuccess(w, resp)
}

// GoogleLogin menangani GET /api/v1/auth/google/login
// Redirect ke Supabase Auth URL untuk login Google
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Bangun URL otorisasi Supabase menggunakan konfigurasi dari constructor
	authURL := fmt.Sprintf("%s/auth/v1/authorize?provider=google&redirect_to=%s",
		h.supabaseURL, h.redirectURL)

	http.Redirect(w, r, authURL, http.StatusFound)
}

// GoogleCallback menangani callback dari Supabase Auth setelah login Google
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Ambil access_token dari query parameter (Supabase redirect)
	tokenStr := r.URL.Query().Get("access_token")
	if tokenStr == "" {
		response.RespondBadRequest(w, "access_token tidak ditemukan di callback")
		return
	}

	// Panggil service untuk login dengan token Supabase
	resp, token, err := h.authService.LoginWithSupabase(r.Context(), tokenStr)
	if err != nil {
		response.RespondUnauthorized(w, err.Error())
		return
	}

	// Token aplikasi dikirim lewat HttpOnly cookie, lalu redirect ke frontend
	h.cookies.Set(w, token, resp.ExpiresIn)

	http.Redirect(w, r, h.frontendURL, http.StatusFound)
}

// Logout menangani POST /api/v1/auth/logout
// Token divalidasi dan dimasukkan ke Redis blacklist agar tidak dapat digunakan lagi
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Ambil token dari HttpOnly cookie
	tokenStr := h.cookies.FromRequest(r)
	if tokenStr == "" {
		response.RespondUnauthorized(w, "cookie autentikasi tidak ditemukan")
		return
	}

	// Cek apakah token sudah di-blacklist sebelumnya
	if h.redisClient != nil {
		alreadyRevoked, err := h.redisClient.IsBlacklisted(r.Context(), tokenStr)
		if err == nil && alreadyRevoked {
			response.RespondUnauthorized(w, "token sudah tidak valid atau telah logout")
			return
		}
	}

	// Validasi token dan ambil claims
	claims, err := h.authService.ValidateTokenWithClaims(tokenStr)
	if err != nil {
		response.RespondUnauthorized(w, err.Error())
		return
	}

	// Hitung sisa waktu token aktif untuk TTL blacklist di Redis
	remainingTTL := time.Until(claims.ExpiresAt.Time)
	if remainingTTL <= 0 {
		remainingTTL = 1 * time.Minute
	}

	// Simpan token ke Redis blacklist
	if h.redisClient != nil {
		if err := h.redisClient.AddToBlacklist(r.Context(), tokenStr, remainingTTL); err != nil {
			response.RespondInternalServerError(w, "gagal memproses logout di blacklist")
			return
		}
	}

	// Hapus HttpOnly cookie autentikasi
	h.cookies.Clear(w)

	// Logout berhasil
	response.RespondSuccess(w, map[string]interface{}{
		"message":    "logout berhasil, token telah dinonaktifkan",
		"user_id":    claims.UserID,
		"revoked_at": time.Now().UTC().Format(time.RFC3339),
	})
}
