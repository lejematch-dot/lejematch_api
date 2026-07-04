package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Environment string
	APIPort     string
	DatabaseURL string
	JWTSecret   string
}

var AppConfigInstance AppConfig

func Load() {

	// Try to load .env file, but don't panic if it's not there
	// (e.g. in Docker, env vars are injected directly)
	envFile := "env"
	_ = godotenv.Load(envFile)

	AppConfigInstance = AppConfig{
		Environment: GetEnv("ENV", "development"),
		APIPort:     GetEnv("API_PORT", "8080"),
		DatabaseURL: GetEnv("DATABASE_URL", ""),
		JWTSecret:   GetEnv("JWT_SECRET", ""),
	}
}

// GetEnv fetches environment variables and returns a fallback value if not found
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
