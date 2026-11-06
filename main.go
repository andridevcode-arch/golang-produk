package main

import (
	"context"
	"log"
	"os"

	"go-produk/internal/app"
)

// main adalah entrypoint untuk Vercel Go Runtime.
// Vercel mendeteksi file main.go di root sebagai entrypoint server,
// dan mewajibkan server listen pada environment variable PORT.
func main() {
	ctx := context.Background()

	appInstance, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := appInstance.StartServer(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
