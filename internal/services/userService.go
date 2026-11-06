package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go-produk/internal/dto"
	"go-produk/internal/models"
	"go-produk/internal/repositories"

	"github.com/jackc/pgx/v5/pgconn"
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type UserService struct{ userRepo *repositories.UserRepository }

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*dto.UserDetailResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if errors.Is(err, repositories.ErrUserNotFound) {
		return nil, errors.New("user tidak ditemukan")
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil profile: %w", err)
	}
	if !user.IsActive {
		return nil, errors.New("akun Anda telah dinonaktifkan, silakan hubungi admin")
	}
	profile, err := s.userRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil user profile: %w", err)
	}
	accounts, err := s.userRepo.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil auth accounts: %w", err)
	}
	return userDetail(user, profile, accounts), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UserDetailResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if errors.Is(err, repositories.ErrUserNotFound) {
		return nil, errors.New("user tidak ditemukan")
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil user: %w", err)
	}
	if !user.IsActive {
		return nil, errors.New("akun Anda telah dinonaktifkan, silakan hubungi admin")
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Username != nil {
		username := strings.ToLower(strings.TrimSpace(*req.Username))
		if len(username) < 3 || len(username) > 50 || !usernamePattern.MatchString(username) {
			return nil, errors.New("username harus 3-50 karakter, hanya huruf kecil, angka, dan tanda hubung")
		}
		user.Username = &username
	}
	if req.Name != "" || req.Username != nil {
		user.UpdatedAt = time.Now()
		var err error
		if req.Username != nil {
			err = s.userRepo.UpdateUsernameAndPublishedPaths(ctx, user)
		} else {
			err = s.userRepo.Update(ctx, user)
		}
		if err != nil {
			if isUniqueViolation(err) {
				return nil, errors.New("username sudah digunakan")
			}
			return nil, fmt.Errorf("gagal memperbarui user: %w", err)
		}
	}
	profile, err := s.userRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil user profile: %w", err)
	}
	if profile == nil {
		profile = &models.UserProfile{UserID: userID}
	}
	if req.DisplayName != "" {
		profile.DisplayName = req.DisplayName
	} else if profile.DisplayName == "" {
		profile.DisplayName = user.Name
	}
	if req.Phone != nil {
		profile.Phone = req.Phone
	}
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = req.AvatarURL
	}
	if req.Location != nil {
		profile.Location = req.Location
	}
	if err := s.userRepo.UpsertProfile(ctx, profile); err != nil {
		return nil, fmt.Errorf("gagal menyimpan profil: %w", err)
	}
	accounts, err := s.userRepo.GetAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil auth accounts: %w", err)
	}
	return userDetail(user, profile, accounts), nil
}

func userDetail(user *models.User, profile *models.UserProfile, accounts []models.AuthAccount) *dto.UserDetailResponse {
	return &dto.UserDetailResponse{ID: user.ID, Name: user.Name, Username: user.Username, Email: user.Email, Role: user.Role, IsActive: user.IsActive, EmailVerifiedAt: user.EmailVerifiedAt, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Profile: profile, Accounts: accounts}
}
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
