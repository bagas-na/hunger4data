package main

import (
	"context"
	"fmt"
	pb "hunger4data/pb/subcription"
	"log"
	"net"
	"os"
	"subscription/internal/adapters/model"
	"subscription/internal/adapters/repo"
	"subscription/internal/service"
	grpcHandler "subscription/internal/transport/grpc"

	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	RedisAddr string
	DBDSN     string
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
		RedisAddr: getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		DBDSN:     getEnv("DB_DSN", ""),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
		GRPCPort:  getEnv("GRPC_PORT", "50052"),
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
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	err = db.AutoMigrate(&model.Subscription{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	db.Exec(`
    CREATE UNIQUE INDEX IF NOT EXISTS idx_user_country_active
    ON subscriptions (user_id, country_code)
    WHERE deleted_at IS NULL
`)

	log.Println("PostgreSQL connected successfully via GORM")
	return db, nil
}

func main() {
	cfg := LoadConfig()

	db, err := NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	service.StartSyncWithoutScheduler(rdb)
	service.StartSyncScheduler(rdb)

	subsRepo := repo.NewSubRepo(db)
	subsServ := service.NewSubService(subsRepo)
	subscriptionService := grpcHandler.NewSubHand(subsServ, rdb)
	grpcServer := grpc.NewServer()

	pb.RegisterSubscription_ServiceServer(grpcServer, &subscriptionService)
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	log.Printf("gRPC server listening on :%s", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
