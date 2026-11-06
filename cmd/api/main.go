package main

import (
	"context"
	"fmt"
	"log"

	"go-produk/internal/app"
)

func main() {
	ctx := context.Background()

	// Load app dan semua koneksi database
	appInstance, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Ambil port dari config, default 8080
	addr := fmt.Sprintf(":%s", appInstance.Config.App.Port)
	if appInstance.Config.App.Port == "" {
		addr = ":8080"
	}

	// Jalankan server
	if err := appInstance.StartServer(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
