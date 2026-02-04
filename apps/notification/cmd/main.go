package main

import (
	notifyv1 "hunger4data/pb/notification"
	"log"
	"net"
	"notification/internal/adapters/db"
	"notification/internal/adapters/mailer"
	"notification/internal/config"
	"notification/internal/service"
	grpcHandler "notification/internal/transport/grpc"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	dbClient, err := db.ConnectDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	mailer := mailer.NewMailersendMailer(cfg.MAILER_SEND_API)
	logRepo := db.NewNotificationLogRepo(dbClient)

	notificationSvc := service.NewNotificationService(mailer, logRepo)

	grpcServer := grpc.NewServer()

	mailerFromName := cfg.MAILER_EMAIL_USERNAME
	mailerFromeEmail := cfg.MAILER_EMAIL_USERNAME + "@" + cfg.MAILER_EMAIL_DOMAIN
	notificationHandler := grpcHandler.NewEmailGRPCServer(notificationSvc, mailerFromName, mailerFromeEmail)

	notifyv1.RegisterEmailServiceServer(grpcServer, notificationHandler)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	log.Printf("payment-service gRPC server listening on %s", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
