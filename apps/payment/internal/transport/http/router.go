package httpHandler

import (
	"payment-service/internal/service"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *service.PaymentService
}

func NewHandler(s *service.PaymentService) *Handler {
	return &Handler{service: s}
}

func RegisterRoutes(e *echo.Echo, paymentSvc *service.PaymentService) {
	h := NewHandler(paymentSvc)

	api := e.Group("/api")

	api.POST("/payments", h.CreatePayment)
	api.GET("/payments/:id", h.GetPayment)
}
