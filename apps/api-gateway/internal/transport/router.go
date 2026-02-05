package http

import (
	"api-gateway/internal/service"

	"github.com/labstack/echo/v4"
)

type SubscriptionHandler interface {
	GetCountries(c echo.Context) error
	GetUserSubs(c echo.Context) error
	CreateSub(c echo.Context) error
	// UpdateSub(c echo.Context) error
	DeleteSub(c echo.Context) error
}

type PaymentHandler interface {
	CreatePayment(c echo.Context) error
	GetPaymentCheckoutURL(c echo.Context) error
	ListPayments(c echo.Context) error
	ListPendingPayments(c echo.Context) error
}

type NotificationHandler interface {
	SendTransactionEmail(c echo.Context) error
}

func AuthRouting(e *echo.Echo, auth *service.AuthHandler) {
	g := e.Group("/auth")
	g.POST("/register", auth.Register) // complete
	g.POST("/login", auth.Login)       // complete
}

func SubscriptionRouting(e *echo.Echo, subscription SubscriptionHandler, secret string) {
	g := e.Group("/subscription")
	g.GET("/countries", subscription.GetCountries) // done
	g.GET("", subscription.GetUserSubs, JWTMiddleware(secret))
	g.POST("", subscription.CreateSub, JWTMiddleware(secret)) // done
	// g.PUT("/:id", subscription.UpdateSub, JWTMiddleware(secret))
	g.DELETE("/:id", subscription.DeleteSub, JWTMiddleware(secret))
}

func PaymentRouting(e *echo.Echo, payment PaymentHandler, secret string) {
	g := e.Group("/payments")
	g.POST("/", payment.CreatePayment)
	g.GET("/checkout/:id", payment.GetPaymentCheckoutURL)
	g.GET("/payments/all", payment.ListPayments)
	g.GET("/payments/pending", payment.ListPendingPayments)
}

func NotificationRouting(e *echo.Echo, notification NotificationHandler, secret string) {
	g := e.Group("/notification")
	g.POST("/", notification.SendTransactionEmail)
}
