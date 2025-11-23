package env

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func GetString(key string, fallback string) string {
	if err := godotenv.Load(); err != nil {
		slog.Warn("the .env file wasn't read -> using default data", "error", err)
	}

	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}