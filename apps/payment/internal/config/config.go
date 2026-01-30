package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort string
	DBDSN    string
}

func Load() *Config {
	env := os.Getenv("ENV")

	// Default to production if ENV is not set
	if env == "" {
		env = "production"
	}

	// Use godotenv only in development
	if env == "development" || env == "dev" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Failed to load .env file. Using existing and/or default ENV values")
		}
	}

	return &Config{
		HttpPort: getEnv("PORT", "9000"),
		DBDSN:    getEnv("DB_DSN", "postgresql://postgres:postgres@localhost:5432/postgres"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Printf("Missing env var: %s. Using fallback value: %s\n", key, fallback)
	}
	return val
}
