package main

import (
	"api-gateway/internal/config"
	"api-gateway/internal/service"
	http "api-gateway/internal/transport"
	authenticatorv1 "hunger4data/pb/authenticator"
	notifyv1 "hunger4data/pb/notification"
	paymentv1 "hunger4data/pb/payment"
	subscriptionv1 "hunger4data/pb/subcription"
	"log"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()
	grpcConnAuth, err := grpc.NewClient(cfg.GRPCAddrAUTH, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to auth gRPC server: %v", err)
	}
	grpcConnNotif, err := grpc.NewClient(cfg.GRPCAddrNOTIFICATION, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to notification gRPC server: %v", err)
	}
	grpcConnPayment, err := grpc.NewClient(cfg.GRPCAddrPAYMENT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to payment gRPC server: %v", err)
	}
	grpcConnSubscription, err := grpc.NewClient(cfg.GRPCAddrSUBSCRIPTION, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to subscription gRPC server: %v", err)
	}
	defer grpcConnAuth.Close()
	defer grpcConnNotif.Close()
	defer grpcConnPayment.Close()
	defer grpcConnSubscription.Close()

	e := echo.New()
	e.HideBanner = true
	e.Use(echomw.RequestLogger())
	e.Use(echomw.Recover())

	authclient := authenticatorv1.NewAuthServiceClient(grpcConnAuth)
	authhand := service.NewHandAuth(authclient)
	http.AuthRouting(e, authhand)

	notifclient := notifyv1.NewEmailServiceClient(grpcConnNotif)
	notifhand := service.NewNotifyHand(notifclient)
	http.NotificationRouting(e, notifhand, cfg.JWTSecret)

	paymentclient := paymentv1.NewPaymentServiceClient(grpcConnPayment)
	paymenthand := service.NewPaymentHand(paymentclient)
	http.PaymentRouting(e, paymenthand, cfg.JWTSecret)

	subscriptionclient := subscriptionv1.NewSubscription_ServiceClient(grpcConnSubscription)
	subscriptionhand := service.NewHandSubs(subscriptionclient)
	http.SubscriptionRouting(e, subscriptionhand, cfg.JWTSecret)

	log.Printf("REST server listening on :%s", cfg.RESTPort)
	if err := e.Start(":" + cfg.RESTPort); err != nil {
		log.Fatalf("failed to start REST server: %v", err)
	}

}
