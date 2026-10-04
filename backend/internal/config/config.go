package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string // HTTP listen port
	DatabaseURL string // Database connection URL
	JWTSecret   string // Secret key for JWT signing
	CORSOrigin  string // Allowed CORS origin
}

// load reads config, fallling back to devlopment defaults. Every value has a default so we can boot the app without wriging a .env first

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://referhub:referhub@localhost:5432/referhub?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-only-insecure-secret"),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3001"),
	}
	if cfg.JWTSecret == "dev-only-insecure-secret" {
		log.Println("WARNING: Using default JWT secret. This is insecure and should be changed in production.")
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
