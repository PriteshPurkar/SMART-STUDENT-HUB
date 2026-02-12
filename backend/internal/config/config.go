package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	FrontendOrigin string
	Environment    string
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Build from parts if DATABASE_URL not set
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "5432")
		dbUser := getEnv("DB_USER", "postgres")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := getEnv("DB_NAME", "scalable_learning")

		if dbPassword == "" {
			dbPassword = "postgres"
		}

		databaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser, dbPassword, dbHost, dbPort, dbName,
		)
	}

	cfg := &Config{
		Port:           getEnv("API_PORT", "8080"),
		DatabaseURL:    databaseURL,
		RedisURL:       os.Getenv("REDIS_URL"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		FrontendOrigin: getEnv("FRONTEND_ORIGIN", "http://localhost:5173"),
		Environment:    getEnv("ENVIRONMENT", "development"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

