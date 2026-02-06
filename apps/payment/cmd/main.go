package main

import (
	"context"
	notifyv1 "hunger4data/pb/notification"
	paymentv1 "hunger4data/pb/payment"
	"log"
	"net"
	"os"
	"os/signal"
	"payment-service/internal/adapters/db"
	"payment-service/internal/adapters/notification"
	stripeAdapter "payment-service/internal/adapters/stripe"
	"payment-service/internal/config"
	"payment-service/internal/service"
	grpcHandler "payment-service/internal/transport/grpc"
	httpHandler "payment-service/internal/transport/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	dbClient, err := db.ConnectDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("db error: %v", err)
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

	stripeAdapter := stripeAdapter.NewStripeAdapter(cfg.STRIPE_SECRET_KEY)

	paymentRepo := db.NewPaymentRepo(dbClient)
	paymentService := service.NewPaymentService(paymentRepo, stripeAdapter, mailer)

	go startGRPCServer(cfg.GRPCPort, paymentService)
	go startHTTPServer(cfg.WebhookPort, paymentService, cfg.STRIPE_WEBHOOK_SECRET)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down webhook and grpc server")
}

func startGRPCServer(grpcPort string, svc service.PaymentService) {
	grpcServer := grpc.NewServer()

	paymentHandler := grpcHandler.NewPaymentGRPCServer(svc)
	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", grpcPort, err)
	}

	log.Printf("payment-service gRPC server listening on %s", grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func startHTTPServer(webhookPort string, svc service.PaymentService, webhookSecret string) {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.POST("/webhooks/stripe", httpHandler.StripeWebhookHandler(webhookSecret, svc))

	log.Println("HTTP listening on", webhookPort)
	log.Fatal(e.Start(":" + webhookPort))
}
