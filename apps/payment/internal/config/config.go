package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort               string
	DBDSN                  string
	STRIPE_SECRET_KEY      string
	STRIPE_PUBLISHABLE_KEY string
	STRIPE_WEBHOOK_SECRET  string
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
		HttpPort:               getEnv("PORT", "9000"),
		DBDSN:                  getEnv("DB_DSN", "postgresql://postgres:postgres@localhost:5432/postgres"),
		STRIPE_SECRET_KEY:      getEnv("STRIPE_SECRET_KEY", ""),
		STRIPE_PUBLISHABLE_KEY: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		STRIPE_WEBHOOK_SECRET:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Printf("Missing env var: %s. Using fallback value: %q\n", key, fallback)
	}
	return val
}
