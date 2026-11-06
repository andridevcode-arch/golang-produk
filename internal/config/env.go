package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// LoadEnv memuat file .env ke environment variables.
// Jika file .env tidak ada, aplikasi tetap lanjut memakai
// environment variables dari sistem.
func LoadEnv() error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		// .env memang opsional, abaikan jika tidak ada
		return nil
	}
	return nil
}

// getEnv mengambil env var, dengan fallback value jika kosong
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return fallback
}

// getEnvAsInt mengambil env var sebagai int
func getEnvAsInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return n
		}
	}
	return fallback
}

// getEnvAsDuration mengambil env var sebagai time.Duration
// contoh value: "10s", "30m", "1h", "24h"
func getEnvAsDuration(key, fallback string) time.Duration {
	v := getEnv(key, fallback)
	// Support "1d" format, convert to "24h"
	v = strings.ReplaceAll(v, "d", "h")
	d, err := time.ParseDuration(v)
	if err != nil {
		f := strings.ReplaceAll(fallback, "d", "h")
		d, _ = time.ParseDuration(f)
	}
	return d
}

// getEnvAsBool mengambil env var sebagai bool
func getEnvAsBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(strings.TrimSpace(v)); err == nil {
			return b
		}
	}
	return fallback
}
