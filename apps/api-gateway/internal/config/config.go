package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCPortAUTH         string
	GRPCPortPAYMENT      string
	GRPCPortSUBSCRIPTION string
	GRPCPortNOTIFICATION string
	RESTPort             string
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
		GRPCPortAUTH:         getEnv("GRPC_PORTAUTH", "50051"),
		GRPCPortNOTIFICATION: getEnv("GRPC_PORT", "50052"),
		GRPCPortPAYMENT:      getEnv("GRPC_PORT", "9000"),
		GRPCPortSUBSCRIPTION: getEnv("GRPC_PORTSUBSCRIPTION", "50054"),
		RESTPort:             getEnv("RESTPort", "8080"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Printf("Missing env var: %s. Using fallback value: %q\n", key, fallback)
	}
	return val
}
