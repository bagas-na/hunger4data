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

	_ "api-gateway/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           Hunger4Data API
// @version         0.0.1
// @description     API for food security data, subscriptions, and donations.
// @description     Provides country-level food security indicators aggregated from HAPI/HDX.
// @description     Users can subscribe to countries, receive daily statistics, and donate to supported causes.
// @termsOfService  https://github.com/bagas-na/hunger4data

// @contact.name    Hunger4Data Backend
// @contact.url     https://github.com/bagas-na/hunger4data
// @contact.email   none

// @license.name  MPL 2.0
// @license.url   https://github.com/bagas-na/hunger4data/blob/main/LICENSE

// @host      localhost:9000
// @BasePath  /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description     Bearer token for authenticated user, subscription, and donation endpoints.
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

	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
