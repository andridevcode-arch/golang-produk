package config

import (
	"time"
)

// Config adalah struct utama yang menampung semua konfigurasi aplikasi
type Config struct {
	App                  AppConfig
	Postgres             PostgresConfig
	Redis                RedisConfig
	Mongo                MongoConfig
	JWTSecret            string
	JWTExpiresIn         time.Duration
	SupabaseURL          string
	SupabaseServiceRoleKey string
	GoogleRedirectURL    string
}

// AppConfig berisi konfigurasi aplikasi umum
type AppConfig struct {
	Environment string
	Port        string
}

// PostgresConfig berisi konfigurasi untuk Supabase PostgreSQL
// Menggunakan single connection string URL
type PostgresConfig struct {
	ConnectionString string
	MaxConns         int32
	MinConns         int32
	MaxConnLifetime  time.Duration
	MaxConnIdleTime  time.Duration
}

// RedisConfig berisi konfigurasi untuk Redis
// Menggunakan single connection string URL
type RedisConfig struct {
	ConnectionString string
	PoolSize         int
}

// MongoConfig berisi konfigurasi untuk MongoDB
type MongoConfig struct {
	URI         string
	Database    string
	MaxPoolSize uint64
	MinPoolSize uint64
	Timeout     time.Duration
}

// Load memuat konfigurasi dari environment variables
func Load() (*Config, error) {
	if err := LoadEnv(); err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("APP_PORT", "8080"),
		},
		Postgres: PostgresConfig{
			ConnectionString: getEnv("DATABASE_URL", ""),
			MaxConns:         int32(getEnvAsInt("POSTGRES_MAX_CONNS", 25)),
			MinConns:         int32(getEnvAsInt("POSTGRES_MIN_CONNS", 5)),
			MaxConnLifetime:  getEnvAsDuration("POSTGRES_MAX_CONN_LIFETIME", "1h"),
			MaxConnIdleTime:  getEnvAsDuration("POSTGRES_MAX_CONN_IDLE_TIME", "30m"),
		},
		Redis: RedisConfig{
			ConnectionString: getEnv("REDIS_URL", ""),
			PoolSize:         getEnvAsInt("REDIS_POOL_SIZE", 10),
		},
		Mongo: MongoConfig{
			URI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database:    getEnv("MONGO_DATABASE", "go_produk"),
			MaxPoolSize: uint64(getEnvAsInt("MONGO_MAX_POOL_SIZE", 100)),
			MinPoolSize: uint64(getEnvAsInt("MONGO_MIN_POOL_SIZE", 10)),
			Timeout:     getEnvAsDuration("MONGO_TIMEOUT", "10s"),
		},
		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTExpiresIn:         getEnvAsDuration("JWT_EXPIRES_IN", "1d"),
		SupabaseURL:          getEnv("SUPABASE_URL", ""),
		SupabaseServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		GoogleRedirectURL:    getEnv("GOOGLE_REDIRECT_URL", ""),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
