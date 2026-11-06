package services

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go-produk/internal/constants"
	"go-produk/internal/dto"
	"go-produk/internal/models"
	"go-produk/internal/repositories"
	"go-produk/pkg/crypto"
	"go-produk/pkg/token"
)

// AuthService menangani business logic autentikasi
type AuthService struct {
	userRepo  *repositories.UserRepository
	tokenMgr  *token.Manager
	jwtExpiry time.Duration
	supabaseURL string
	supabaseServiceRoleKey string
}

// NewAuthService membuat instance AuthService baru
func NewAuthService(userRepo *repositories.UserRepository, tokenMgr *token.Manager, jwtExpiry time.Duration, supabaseURL, supabaseServiceRoleKey string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenMgr:  tokenMgr,
		jwtExpiry: jwtExpiry,
		supabaseURL: supabaseURL,
		supabaseServiceRoleKey: supabaseServiceRoleKey,
	}
}

// generateUUID membuat UUID v4 sederhana menggunakan crypto/rand
func generateUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("unknown-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// Register mendaftarkan user baru
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Validasi input
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("name, email, dan password wajib diisi")
	}
	if len(req.Password) < 6 {
		return nil, fmt.Errorf("password minimal 6 karakter")
	}

	// Cek email sudah terdaftar
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("email sudah terdaftar")
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("gagal menghash password: %w", err)
	}

	// Buat user baru
	now := time.Now()
	user := &models.User{
		ID:           generateUUID(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         constants.RoleUser,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Simpan ke database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal menyimpan user: %w", err)
	}

	// Catat otomatis auth_account dengan provider email
	_ = s.userRepo.CreateAccount(ctx, &models.AuthAccount{
		UserID:            user.ID,
		Provider:          "email",
		ProviderAccountID: user.Email,
		Email:             user.Email,
	})

	// Generate JWT token
	tokenStr, err := s.tokenMgr.Generate(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %w", err)
	}

	return &dto.LoginResponse{
		User:        *user,
		AccessToken: tokenStr,
		ExpiresIn:   int(s.jwtExpiry.Seconds()),
	}, nil
}

// ValidateToken validates JWT token and returns user ID
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	claims, err := s.tokenMgr.Validate(tokenString)
	if err != nil {
		return "", fmt.Errorf("token tidak valid atau sudah kadaluarsa")
	}
	return claims.UserID, nil
}

// ValidateTokenWithClaims validates JWT token and returns full claims
func (s *AuthService) ValidateTokenWithClaims(tokenString string) (*token.Claims, error) {
	claims, err := s.tokenMgr.Validate(tokenString)
	if err != nil {
		return nil, fmt.Errorf("token tidak valid atau sudah kadaluarsa")
	}
	return claims, nil
}

// Login melakukan autentikasi user
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Validasi input
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email dan password wajib diisi")
	}

	// Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, fmt.Errorf("email tidak terdaftar")
		}
		return nil, fmt.Errorf("terjadi kesalahan pada sistem")
	}

	// Verifikasi password
	if !crypto.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("password salah")
	}

	// Cek status aktif akun
	if !user.IsActive {
		return nil, fmt.Errorf("akun Anda telah dinonaktifkan, silakan hubungi admin")
	}

	// Generate JWT token
	tokenStr, err := s.tokenMgr.Generate(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %w", err)
	}

	return &dto.LoginResponse{
		User:        *user,
		AccessToken: tokenStr,
		ExpiresIn:   int(s.jwtExpiry.Seconds()),
	}, nil
}

// supabaseUser adalah response dari Supabase Auth /auth/v1/user
type supabaseUser struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	EmailConfirmedAt *time.Time `json:"email_confirmed_at"`
	Phone         string     `json:"phone"`
	RawUserMetaData struct {
		FullName string `json:"full_name"`
		AvatarURL string `json:"avatar_url"`
	} `json:"raw_user_meta_data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginWithSupabase melakukan autentikasi user melalui token Supabase Auth (Google OAuth)
func (s *AuthService) LoginWithSupabase(ctx context.Context, supabaseAccessToken string) (*dto.LoginResponse, error) {
	if supabaseAccessToken == "" {
		return nil, fmt.Errorf("supabase access token wajib diisi")
	}

	// Panggil Supabase Auth API untuk mendapatkan profil user
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.supabaseURL+"/auth/v1/user", nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request ke Supabase: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+supabaseAccessToken)
	req.Header.Set("apikey", s.supabaseServiceRoleKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal memanggil Supabase Auth API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Supabase Auth gagal (status %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var su supabaseUser
	if err := json.NewDecoder(resp.Body).Decode(&su); err != nil {
		return nil, fmt.Errorf("gagal decode response Supabase: %w", err)
	}

	if su.Email == "" {
		return nil, fmt.Errorf("email tidak tersedia dari Supabase Auth")
	}

	// Cari user berdasarkan email di database lokal
	user, err := s.userRepo.FindByEmail(ctx, su.Email)
	if err != nil {
		if !errors.Is(err, repositories.ErrUserNotFound) {
			return nil, fmt.Errorf("terjadi kesalahan pada sistem")
		}

		// Buat user baru jika belum ada
		now := time.Now()
		user = &models.User{
			ID:           generateUUID(),
			Name:         su.RawUserMetaData.FullName,
			Email:        su.Email,
			PasswordHash: "",
			Role:         constants.RoleUser,
			IsActive:     true,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("gagal menyimpan user: %w", err)
		}

		// Catat auth_account dengan provider google
		_ = s.userRepo.CreateAccount(ctx, &models.AuthAccount{
			UserID:            user.ID,
			Provider:          "google",
			ProviderAccountID: su.ID,
			Email:             su.Email,
			AvatarURL:         &su.RawUserMetaData.AvatarURL,
		})
	} else {
		// Pastikan auth_account google sudah tercatat
		accounts, err := s.userRepo.GetAccountsByUserID(ctx, user.ID)
		if err == nil {
			hasGoogle := false
			for _, acc := range accounts {
				if acc.Provider == "google" && acc.ProviderAccountID == su.ID {
					hasGoogle = true
					break
				}
			}
			if !hasGoogle {
				_ = s.userRepo.CreateAccount(ctx, &models.AuthAccount{
					UserID:            user.ID,
					Provider:          "google",
					ProviderAccountID: su.ID,
					Email:             su.Email,
					AvatarURL:         &su.RawUserMetaData.AvatarURL,
				})
			}
		}
	}

	// Generate JWT token aplikasi
	tokenStr, err := s.tokenMgr.Generate(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("gagal generate token: %w", err)
	}

	return &dto.LoginResponse{
		User:        *user,
		AccessToken: tokenStr,
		ExpiresIn:   int(s.jwtExpiry.Seconds()),
	}, nil
}
