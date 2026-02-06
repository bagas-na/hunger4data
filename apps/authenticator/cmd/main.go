package main

import (
	"authenticator/internal/adapters/crypto"
	"authenticator/internal/adapters/notification"
	"authenticator/internal/adapters/repo"
	"authenticator/internal/service"
	handler "authenticator/internal/transport/grpc"
	"fmt"
	authenticatorv1 "hunger4data/pb/authenticator"
	notifyv1 "hunger4data/pb/notification"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBDSN                string
	JWTSecret            string
	JWTDuration          time.Duration
	GRPCPort             string
	NotificationGRPCAddr string
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

	durString := getEnv("JWT_DURATION", "60m")

	duration, err := time.ParseDuration(durString)
	if err != nil {
		fmt.Printf("%q is not a valid duration. Using default value of 60m", durString)
		duration = 60 * time.Minute
	}

	return &Config{
		GRPCPort:             getEnv("GRPC_PORT", "50051"),
		DBDSN:                getEnv("DB_DSN", ""),
		JWTSecret:            getEnv("JWT_SECRET", "change-this-super-secret-but-insecure-key"),
		JWTDuration:          duration,
		NotificationGRPCAddr: getEnv("GRPC_ADDR_NOTIFICATION", ":50052"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func NewDBConnection(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&repo.User{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	fmt.Println("user automigration complete")

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

	grpcConnNotif, err := grpc.NewClient(
		cfg.NotificationGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to notification gRPC server: %v", err)
	}
	defer grpcConnNotif.Close()

	notifclient := notifyv1.NewEmailServiceClient(grpcConnNotif)
	mailer := notification.NewMailer(notifclient)

	cf := crypto.NewJwtPass(cfg.JWTSecret, cfg.JWTDuration)

	userRepo := repo.NewUserRepo(db)
	authserv := service.NewAuthService(userRepo, cf, mailer)
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
