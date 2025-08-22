package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		DatabaseURL: fmt.Sprintf(
			"postgres://%s:%s@%s:5432/%s?sslmode=disable",
			getEnv("DB_USER", "user"),
			getEnv("DB_PASSWORD", "password"),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_NAME", "spy_cat_agency"),
		),
		Port: getEnv("PORT", "8080"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		return strconv.Atoi(value)
	}
	return defaultValue, nil
}
