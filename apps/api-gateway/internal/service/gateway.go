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
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	req := &paymentv1.CreatePaymentRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	req.UserId = c.Get("user_id").(string)

	resp, err := h.paymentclient.CreatePayment(ctx, req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, CreatePaymentResponse{
		Message: "Payment created successfully.",
		PaymentInfo: PaymentInfo{
			PaymentID:       resp.Payment.Id,
			CountryCode:     resp.Payment.CountryCode,
			TransactionType: resp.Payment.TransactionType,
			Amount:          resp.Payment.Amount,
			Currency:        resp.Payment.Currency,
			Status:          resp.Payment.Status,
			CreatedAt:       resp.Payment.CreatedAt,
			UpdatedAt:       resp.Payment.UpdatedAt,
		},
		CheckoutUrl: resp.CheckoutUrl,
	})
}

func (h *PaymentHandler) GetPaymentCheckoutURL(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	req := &paymentv1.GetPaymentCheckoutURLRequest{
		PaymentId: id,
	}
	resp, err := h.paymentclient.GetPaymentCheckoutURL(ctx, req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) ListPayments(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	req := &paymentv1.ListPaymentsRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	req.UserId = c.Get("user_id").(string)

	resp, err := h.paymentclient.ListPayments(ctx, req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	var paymentList []PaymentInfo

	for _, p := range resp.Payments {
		paymentInfo := PaymentInfo{
			PaymentID:       p.Id,
			CountryCode:     p.CountryCode,
			TransactionType: p.TransactionType,
			Amount:          p.Amount,
			Currency:        p.Currency,
			Status:          p.Status,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
		}

		paymentList = append(paymentList, paymentInfo)
	}

	return c.JSON(http.StatusOK, paymentList)
}

func (h *PaymentHandler) ListPendingPayments(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	req := &paymentv1.ListPendingPaymentsRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	req.UserId = c.Get("user_id").(string)

	resp, err := h.paymentclient.ListPendingPayments(ctx, req)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	var paymentList []PaymentInfo

	for _, p := range resp.Payments {
		paymentInfo := PaymentInfo{
			PaymentID:       p.Id,
			CountryCode:     p.CountryCode,
			TransactionType: p.TransactionType,
			Amount:          p.Amount,
			Currency:        p.Currency,
			Status:          p.Status,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
		}

		paymentList = append(paymentList, paymentInfo)
	}

	return c.JSON(http.StatusOK, paymentList)
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
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	req := &notifyv1.SendEmailRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	resp, err := h.notificationclient.SendTransactionEmail(ctx, req)
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
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	req := &authenticatorv1.LoginRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

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
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	resp, err := h.subscriptionclient.Get_Countries(ctx, &pb.Empty{})
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": resp.Countries,
	})
}

func (h *subscriptionHandler) CreateSub(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	userId := c.Get("user_id").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "userId is required"})
	}

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

// func (h *subscriptionHandler) UpdateSub(c echo.Context) error {
// 	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
// 	defer cancel()

// 	id := c.Param("id")
// 	if id == "" {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
// 	}
// 	req := &pb.Subscription_Request{}
// 	if err := c.Bind(req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{
// 			"error": "invalid request body",
// 		})
// 	}

// 	resp, err := h.subscriptionclient.Update_Subscription(ctx, &pb.Subscription_Request{
// 		Id:          id,
// 		UserId:      req.UserId,
// 		CountryCode: req.CountryCode,
// 	})
// 	if err != nil {
// 		return utils.MapGRPCError(c, err)
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"message": resp.Message,
// 	})
// }

func (h *subscriptionHandler) DeleteSub(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	subriptionId := c.Param("id")

	userId := c.Get("user_id").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "userId is required"})
	}

	if subriptionId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "subcriptionId is required"})
	}

	resp, err := h.subscriptionclient.Delete_Subscription(ctx, &pb.Subscription_Request{
		Id:     subriptionId,
		UserId: userId,
	})
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": resp.Message,
	})
}
