package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-produk/internal/database"
	"go-produk/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct{ pool *pgxpool.Pool }

func NewUserRepository(pool *pgxpool.Pool) *UserRepository { return &UserRepository{pool: pool} }

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, name, email, password_hash, role, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := r.pool.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func scanUser(row pgx.Row) (*models.User, error) {
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.IsDeleted, &user.EmailVerifiedAt, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := scanUser(r.pool.QueryRow(ctx, `SELECT id,name,username,email,password_hash,role,is_active,is_deleted,email_verified_at,created_at,updated_at FROM users WHERE email=$1 AND is_deleted=false`, email))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	user, err := scanUser(r.pool.QueryRow(ctx, `SELECT id,name,username,email,password_hash,role,is_active,is_deleted,email_verified_at,created_at,updated_at FROM users WHERE id=$1 AND is_deleted=false`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET name=$2, username=$3, email=$4, role=$5, updated_at=now() WHERE id=$1`, user.ID, user.Name, user.Username, user.Email, user.Role)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdateUsernameAndPublishedPaths(ctx context.Context, user *models.User) error {
	return database.WithinTx(ctx, r.pool, func(ctx context.Context, tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE users SET name=$2,username=$3,email=$4,role=$5,updated_at=now() WHERE id=$1`, user.ID, user.Name, user.Username, user.Email, user.Role); err != nil {
			return fmt.Errorf("update user: %w", err)
		}
		if _, err := tx.Exec(ctx, `UPDATE published_pages SET custom_path='/' || $2 || '/' || slug,updated_at=now() WHERE user_id=$1`, user.ID, *user.Username); err != nil {
			return fmt.Errorf("update published paths: %w", err)
		}
		return nil
	})
}

func (r *UserRepository) GetProfileByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
	profile := &models.UserProfile{}
	err := r.pool.QueryRow(ctx, `SELECT id,user_id,display_name,phone,bio,avatar_url,location,created_at,updated_at FROM user_profiles WHERE user_id=$1`, userID).Scan(&profile.ID, &profile.UserID, &profile.DisplayName, &profile.Phone, &profile.Bio, &profile.AvatarURL, &profile.Location, &profile.CreatedAt, &profile.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user profile: %w", err)
	}
	return profile, nil
}

func (r *UserRepository) UpsertProfile(ctx context.Context, profile *models.UserProfile) error {
	return r.pool.QueryRow(ctx, `INSERT INTO user_profiles (user_id,display_name,phone,bio,avatar_url,location,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,now(),now()) ON CONFLICT (user_id) DO UPDATE SET display_name=EXCLUDED.display_name,phone=EXCLUDED.phone,bio=EXCLUDED.bio,avatar_url=EXCLUDED.avatar_url,location=EXCLUDED.location,updated_at=now() RETURNING id,created_at,updated_at`, profile.UserID, profile.DisplayName, profile.Phone, profile.Bio, profile.AvatarURL, profile.Location).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)
}

func (r *UserRepository) GetAccountsByUserID(ctx context.Context, userID string) ([]models.AuthAccount, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,user_id,provider,provider_account_id,email,avatar_url,created_at,updated_at FROM auth_accounts WHERE user_id=$1 ORDER BY created_at ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("get auth accounts: %w", err)
	}
	defer rows.Close()
	accounts := []models.AuthAccount{}
	for rows.Next() {
		var a models.AuthAccount
		if err := rows.Scan(&a.ID, &a.UserID, &a.Provider, &a.ProviderAccountID, &a.Email, &a.AvatarURL, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan auth account: %w", err)
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *UserRepository) CreateAccount(ctx context.Context, acc *models.AuthAccount) error {
	err := r.pool.QueryRow(ctx, `INSERT INTO auth_accounts (user_id,provider,provider_account_id,email,avatar_url,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,now(),now()) ON CONFLICT (provider,provider_account_id) DO NOTHING RETURNING id,created_at,updated_at`, acc.UserID, acc.Provider, acc.ProviderAccountID, acc.Email, acc.AvatarURL).Scan(&acc.ID, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("create auth account: %w", err)
	}
	return nil
}

func (r *UserRepository) SoftDelete(ctx context.Context, id string) error {
	result, err := r.pool.Exec(ctx, `UPDATE users SET is_deleted=true,updated_at=now() WHERE id=$1 AND is_deleted=false`, id)
	if err != nil {
		return fmt.Errorf("soft delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}
