package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPCAddrAUTH         string
	GRPCAddrPAYMENT      string
	GRPCAddrSUBSCRIPTION string
	GRPCAddrNOTIFICATION string
	RESTPort             string
	JWTSecret            string
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
		GRPCAddrAUTH:         getEnv("GRPC_ADDR_AUTH", ":50051"),
		GRPCAddrNOTIFICATION: getEnv("GRPC_ADDR_NOTIFICATION", ":50052"),
		GRPCAddrPAYMENT:      getEnv("GRPC_ADDR_PAYMENT", ":50053"),
		GRPCAddrSUBSCRIPTION: getEnv("GRPC_ADDR_SUBSCRIPTION", ":50054"),
		RESTPort:             getEnv("RESTPort", "8080"),
		JWTSecret:            getEnv("JWT_SECRET", "change-this-super-secret-but-insecure-key"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Printf("Missing env var: %s. Using fallback value: %q\n", key, fallback)
		return fallback
	}
	return val
}
