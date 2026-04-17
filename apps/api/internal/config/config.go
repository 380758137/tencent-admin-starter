package config

import (
	"errors"
	"os"
)

type Config struct {
	Port         string
	MySQLDSN     string
	JWTSecret    string
	AppEnv       string
	MigrationDir string
}

func Load() (Config, error) {
	cfg := Config{
		Port:         envOr("APP_PORT", "8080"),
		MySQLDSN:     envOr("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/tencent_admin?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:    envOr("JWT_SECRET", ""),
		AppEnv:       envOr("APP_ENV", "development"),
		MigrationDir: envOr("MIGRATION_DIR", "migrations"),
	}
	if cfg.JWTSecret == "" || cfg.JWTSecret == "change-me-in-production" {
		return Config{}, errors.New("JWT_SECRET must be explicitly configured")
	}
	return cfg, nil
}

func envOr(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
