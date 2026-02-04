package main

import (
	"authenticator/internal/adapters/crypto"
	"authenticator/internal/adapters/repo"
	"authenticator/internal/service"
	handler "authenticator/internal/transport/grpc"
	"fmt"
	authenticatorv1 "hunger4data/pb/authenticator"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBDSN string
	// DBHost     string
	// DBPort     string
	// DBUser     string
	// DBPassword string
	// DBName     string
	JWTSecret string
	GRPCPort  string
	RESTPort  string
}

func LoadConfig() *Config {
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
		DBDSN: getEnv("DB_DSN", ""),
		// DBHost:     getEnv("DB_HOST", "localhost"),
		// DBPort:     getEnv("DB_PORT", "5432"),
		// DBUser:     getEnv("DB_USER", "postgres"),
		// DBPassword: getEnv("DB_PASSWORD", "1"),
		// DBName:     getEnv("DB_NAME", "test"),
		JWTSecret: getEnv("JWT_SECRET", "change-this-super-secret-but-insecure-key"),
		GRPCPort:  getEnv("GRPC_PORT", "50051"),
		RESTPort:  getEnv("REST_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func NewDBConnection(cfg *Config) (*gorm.DB, error) {
	fmt.Printf("Value of cfg:\n%#v", cfg)
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
	// 	cfg.DBHost,
	// 	cfg.DBUser,
	// 	cfg.DBPassword,
	// 	cfg.DBName,
	// 	cfg.DBPort,
	// )
	dsn := cfg.DBDSN

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// // 3. Access the underlying generic sql.DB to set pool settings
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	// }

	// sqlDB.SetMaxOpenConns(25)
	// sqlDB.SetMaxIdleConns(25)
	// sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	log.Println("PostgreSQL connected successfully via GORM")
	return db, nil
}

func main() {
	cfg := LoadConfig()

	db, err := NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&repo.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	fmt.Println("user automigration complete")

	userRepo := repo.NewUserRepo(db)

	cf := crypto.NewJwtPass(cfg.JWTSecret)
	authserv := service.NewAuthService(userRepo, cf)
	authhand := handler.NewHandService(authserv)
	grpcServer := grpc.NewServer()

	authenticatorv1.RegisterAuthServiceServer(grpcServer, &authhand)
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	log.Printf("gRPC server listening on :%s", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
