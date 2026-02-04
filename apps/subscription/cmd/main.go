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
	DBDSN     string
	JWTSecret string
	GRPCPort  string
	RESTPort  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
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
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	service.StartSyncScheduler(rdb)
	err = db.AutoMigrate(&model.Subscription{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
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
