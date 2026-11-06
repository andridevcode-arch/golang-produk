package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-produk/internal/config"
	"go-produk/internal/constants"
	"go-produk/internal/models"
	"go-produk/internal/repositories"
	"go-produk/pkg/crypto"
)

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("muat konfigurasi: %v", err)
	}

	// Connect ke PostgreSQL pool
	pool, err := pgxpool.New(ctx, cfg.Postgres.ConnectionString)
	if err != nil {
		log.Fatalf("create connection pool: %v", err)
	}
	defer pool.Close()

	userRepo := repositories.NewUserRepository(pool)

	// Cek apakah user admin sudah ada
	existingUser, err := userRepo.FindByEmail(ctx, "admin@go-produk.local")
	if err != nil && !errors.Is(err, repositories.ErrUserNotFound) {
		log.Fatalf("gagal mencari user: %v", err)
	}

	// Hash password admin
	password := "Admin@123"
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		log.Fatalf("gagal hash password: %v", err)
	}

	// Upsert user admin
	var userID string
	if existingUser != nil {
		// Update user yang sudah ada
		existingUser.PasswordHash = hashedPassword
		existingUser.Role = constants.RoleAdmin
		existingUser.IsActive = true
		existingUser.UpdatedAt = time.Now()
		if err := userRepo.Update(ctx, existingUser); err != nil {
			log.Fatalf("gagal update user: %v", err)
		}
		userID = existingUser.ID
	} else {
		// Buat user baru
		user := &models.User{
			ID:           generateUUID(),
			Name:         "Admin",
			Email:        "admin@go-produk.local",
			PasswordHash: hashedPassword,
			Role:         constants.RoleAdmin,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatalf("gagal create user: %v", err)
		}
		userID = user.ID
	}

	// Simpan auth_account
	if err := userRepo.CreateAccount(ctx, &models.AuthAccount{
		UserID:            userID,
		Provider:          "email",
		ProviderAccountID: "admin@go-produk.local",
		Email:             "admin@go-produk.local",
	}); err != nil {
		log.Fatalf("gagal create auth_account: %v", err)
	}

	// Print kredensial login
	fmt.Println("Admin seeder berhasil!")
	fmt.Println("Email: admin@go-produk.local")
	fmt.Println("Password: Admin@123")
}
