package config

import (
	"errors"
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
