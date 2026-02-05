package service

import (
	"api-gateway/internal/utils"
	"context"
	authenticatorv1 "hunger4data/pb/authenticator"
	notifyv1 "hunger4data/pb/notification"
	paymentv1 "hunger4data/pb/payment"
	pb "hunger4data/pb/subcription"
	subscriptionv1 "hunger4data/pb/subcription"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	paymentclient paymentv1.PaymentServiceClient
}

func NewPaymentHand(paymentclient paymentv1.PaymentServiceClient) *PaymentHandler {
	return &PaymentHandler{
		paymentclient: paymentclient,
	}
}

func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	req := &paymentv1.CreatePaymentRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	resp, err := h.paymentclient.CreatePayment(c.Request().Context(), req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) GetPaymentCheckoutURL(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	req := &paymentv1.GetPaymentCheckoutURLRequest{
		PaymentId: id,
	}
	resp, err := h.paymentclient.GetPaymentCheckoutURL(c.Request().Context(), req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) ListPayments(c echo.Context) error {
	req := &paymentv1.ListPaymentsRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	resp, err := h.paymentclient.ListPayments(c.Request().Context(), req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) ListPendingPayments(c echo.Context) error {
	req := &paymentv1.ListPendingPaymentsRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	resp, err := h.paymentclient.ListPendingPayments(c.Request().Context(), req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

type NotificationHandler struct {
	notificationclient notifyv1.EmailServiceClient
}

func NewNotifyHand(notificationclient notifyv1.EmailServiceClient) *NotificationHandler {
	return &NotificationHandler{
		notificationclient: notificationclient,
	}
}

func (h *NotificationHandler) SendTransactionEmail(c echo.Context) error {
	req := &notifyv1.SendEmailRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	resp, err := h.notificationclient.SendTransactionEmail(c.Request().Context(), req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

type AuthHandler struct {
	authclient authenticatorv1.AuthServiceClient
}

func NewHandAuth(authclient authenticatorv1.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		authclient: authclient,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := &authenticatorv1.LoginRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	ctx := context.TODO()

	resp, err := h.authclient.Login(ctx, req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": resp.Token,
	})
}

func (h *AuthHandler) Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	req := &authenticatorv1.RegisterRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing username or password",
		})
	}

	resp, err := h.authclient.Register(ctx, req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": resp.Message,
		"user": map[string]interface{}{
			"id":       resp.User.Id,
			"username": resp.User.Username,
		},
	})
}

type subscriptionHandler struct {
	subscriptionclient subscriptionv1.Subscription_ServiceClient
}

func NewHandSubs(subscriptionclient subscriptionv1.Subscription_ServiceClient) *subscriptionHandler {
	return &subscriptionHandler{
		subscriptionclient: subscriptionclient,
	}
}

func (h *subscriptionHandler) GetCountries(c echo.Context) error {

	ctx := context.TODO()
	resp, err := h.subscriptionclient.Get_Countries(ctx, &pb.Empty{})
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": resp.Countries,
	})
}

func (h *subscriptionHandler) CreateSub(c echo.Context) error {
	ctx := context.TODO()

	userId := c.Get("user_id").(string)

	req := &pb.Subscription_Request{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	req.UserId = userId

	resp, err := h.subscriptionclient.Create_Subscription(ctx, req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": resp.Message,
	})
}

func (h *subscriptionHandler) GetUserSubs(c echo.Context) error {
	userId := c.Get("user_id").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	ctx := context.TODO()
	resp, err := h.subscriptionclient.Get_Subscriptions(ctx, &pb.Subscription_Request{
		UserId: userId,
	})

	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": resp.Subscription,
	})
}

func (h *subscriptionHandler) UpdateSub(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	req := &pb.Subscription_Request{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}
	ctx := context.TODO()
	resp, err := h.subscriptionclient.Update_Subscription(ctx, &pb.Subscription_Request{
		Id:          id,
		UserId:      req.UserId,
		CountryCode: req.CountryCode,
	})
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": resp.Message,
	})
}

func (h *subscriptionHandler) DeleteSub(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	ctx := context.TODO()
	resp, err := h.subscriptionclient.Delete_Subscription(ctx, &pb.Subscription_Request{
		Id: id,
	})
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": resp.Message,
	})
}
