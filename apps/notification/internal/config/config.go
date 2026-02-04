package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort              string
	DBDSN                 string
	MAILER_SEND_API       string
	MAILER_EMAIL_USERNAME string
	MAILER_EMAIL_DOMAIN   string
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
		GRPCPort:              getEnv("GRPC_PORT", "50054"),
		DBDSN:                 getEnv("DB_DSN", "postgresql://postgres:postgres@localhost:5432/postgres"),
		MAILER_SEND_API:       getEnv("MAILER_SEND_API_KEY", ""),
		MAILER_EMAIL_DOMAIN:   getEnv("MAILER_EMAIL_DOMAIN", ""),
		MAILER_EMAIL_USERNAME: getEnv("MAILER_EMAIL_USERNAME", ""),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Printf("Missing env var: %s. Using fallback value: %q\n", key, fallback)
	}
	return val
}
