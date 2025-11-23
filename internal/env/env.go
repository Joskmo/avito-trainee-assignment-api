// Package env provides helper functions for reading environment variables.
package env

import (
	"os"
)

// GetString returns the value of the environment variable named by the key.
// If the variable is not present, it returns the fallback value.
func GetString(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
