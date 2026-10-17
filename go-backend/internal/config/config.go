package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// DefaultPort is used when PORT is not set. It matches .env.example and openapi.yaml.
const DefaultPort = "4500"

// Config holds the values the server needs at startup.
type Config struct {
	Port        string
	DatabaseURL string
}

// Load reads PORT and DATABASE_URL from the environment and validates them.
func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return Config{}, fmt.Errorf("PORT must be a number from 1 to 65535, got %q", port)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is not set")
	}

	return Config{Port: port, DatabaseURL: databaseURL}, nil
}
