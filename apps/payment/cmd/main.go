package main

import (
	paymentv1 "hunger4data/pb/payment"
	"log"
	"net"
	"payment-service/internal/adapters/db"
	stripeAdapter "payment-service/internal/adapters/stripe"
	"payment-service/internal/config"
	"payment-service/internal/service"
	grpcHandler "payment-service/internal/transport/grpc"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	dbClient, err := db.ConnectDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}

	stripeAdapter := stripeAdapter.NewStripeAdapter(cfg.STRIPE_SECRET_KEY)
	paymentRepo := db.NewPaymentRepo(dbClient)
	paymentService := service.NewPaymentService(paymentRepo, stripeAdapter)
	paymentHandler := grpcHandler.NewPaymentGRPCServer(paymentService)

	grpcServer := grpc.NewServer()

	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	log.Printf("payment-service gRPC server listening on %s", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
