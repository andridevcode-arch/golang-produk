package config

import (
	"errors"
	"net/http"
	"strings"
)

// Validate memeriksa validitas konfigurasi
func (c *Config) Validate() error {
	if strings.TrimSpace(c.App.Port) == "" {
		return errors.New("APP_PORT wajib diisi")
	}

	// Postgres optional
	if strings.TrimSpace(c.Postgres.ConnectionString) != "" {
		if err := c.validatePostgres(); err != nil {
			return err
		}
	}

	// Redis optional
	if strings.TrimSpace(c.Redis.ConnectionString) != "" {
		if err := c.validateRedis(); err != nil {
			return err
		}
	}

	// JWT config validation
	if strings.TrimSpace(c.JWTSecret) == "" {
		return errors.New("JWT_SECRET wajib diisi")
	}

	// Cookie config validation
	if c.Cookie.SameSite == http.SameSiteNoneMode && !c.Cookie.Secure {
		return errors.New("COOKIE_SAME_SITE=none mewajibkan COOKIE_SECURE=true")
	}

	// CORS config validation
	for _, origin := range c.AllowedOrigins {
		if origin == "*" {
			return errors.New("CORS_ALLOWED_ORIGINS tidak boleh memakai wildcard karena kredensial diaktifkan")
		}
	}

	return nil
}

func (c *Config) validatePostgres() error {
	if strings.TrimSpace(c.Postgres.ConnectionString) == "" {
		return errors.New("DATABASE_URL wajib diisi")
	}
	return nil
}

func (c *Config) validateRedis() error {
	if strings.TrimSpace(c.Redis.ConnectionString) == "" {
		return errors.New("REDIS_URL wajib diisi")
	}
	return nil
}
