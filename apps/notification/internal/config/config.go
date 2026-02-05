package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPort        string
	DBDSN           string
	MAILER_HOST     string
	MAILER_PORT     string
	MAILER_LOGIN    string
	MAILER_PASSWORD string
	MAILER_USERNAME string
	MAILER_EMAIL    string
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
		GRPCPort:        getEnv("GRPC_PORT", "50054"),
		DBDSN:           getEnv("DB_DSN", "postgresql://postgres:postgres@localhost:5432/postgres"),
		MAILER_HOST:     getEnv("MAILER_HOST", ""),
		MAILER_PORT:     getEnv("MAILER_PORT", ""),
		MAILER_LOGIN:    getEnv("MAILER_LOGIN", ""),
		MAILER_PASSWORD: getEnv("MAILER_PASSWORD", ""),
		MAILER_USERNAME: getEnv("MAILER_USERNAME", ""),
		MAILER_EMAIL:    getEnv("MAILER_EMAIL", ""),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Printf("Missing env var: %s. Using fallback value: %q\n", key, fallback)
	}
	return val
}
