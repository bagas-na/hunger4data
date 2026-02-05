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

// CreatePayment godoc
// @Summary      Create a donation payment
// @Description  Creates a payment for donating to a country and returns a checkout URL.
// @Tags         Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body CreatePaymentRequest true "Create payment request"
// @Success      201 {object} CreatePaymentResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /payments [post]
func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	userId := c.Get("user_id").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "userId is required"})
	}

	var req CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	resp, err := h.paymentclient.CreatePayment(ctx,
		&paymentv1.CreatePaymentRequest{
			UserId:      userId,
			CountryCode: req.CountryCode,
			Amount:      req.Amount,
			Currency:    req.Currency,
		})
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusCreated, CreatePaymentResponse{
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

// GetPaymentCheckoutURL godoc
// @Summary      Get checkout URL for a payment
// @Description  Retrieves the checkout URL for an existing payment.
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Payment ID (UUID)"
// @Success      200 {object} GetPaymentCheckoutURLResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /payments/{id}/checkout [get]
func (h *PaymentHandler) GetPaymentCheckoutURL(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	resp, err := h.paymentclient.GetPaymentCheckoutURL(ctx,
		&paymentv1.GetPaymentCheckoutURLRequest{
			PaymentId: id,
		})
	if err != nil {
		return utils.MapGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

// ListPayments godoc
// @Summary      List user payments
// @Description  Returns all payments created by the authenticated user.
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} PaymentInfo
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /payments [get]
func (h *PaymentHandler) ListPayments(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	userId := c.Get("user_id").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "userId is required"})
	}

	resp, err := h.paymentclient.ListPayments(ctx,
		&paymentv1.ListPaymentsRequest{
			UserId: userId,
		},
	)
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

// ListPendingPayments godoc
// @Summary      List pending payments
// @Description  Returns all pending payments for the authenticated user.
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} PaymentInfo
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /payments/pending [get]
func (h *PaymentHandler) ListPendingPayments(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	userId := c.Get("user_id").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "userId is required"})
	}

	resp, err := h.paymentclient.ListPendingPayments(ctx,
		&paymentv1.ListPendingPaymentsRequest{
			UserId: userId,
		},
	)
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
	if err := c.Bind(&req); err != nil {
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

// Login godoc
// @Summary      User login
// @Description  Authenticates a user and returns a JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login request"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	resp, err := h.authclient.Login(ctx,
		&authenticatorv1.LoginRequest{
			Username: req.Username,
			Password: req.Password,
		},
	)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Token: resp.Token,
	})
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Register request"
// @Success      201 {object} RegisterResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	var req RegisterRequest
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

	resp, err := h.authclient.Register(ctx,
		&authenticatorv1.RegisterRequest{
			Username: req.Username,
			Password: req.Password,
		},
	)
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	// return c.JSON(http.StatusCreated, map[string]interface{}{
	// 	"message": resp.Message,
	// 	"user": map[string]interface{}{
	// 		"id":       resp.User.Id,
	// 		"username": resp.User.Username,
	// 	},
	// })

	return c.JSON(http.StatusCreated,
		RegisterResponse{
			Message: resp.Message,
			User: struct {
				Id       string `json:"id" example:"45dcf27b-29e4-46aa-a3aa-6e2b33ca86ae"`
				Username string `json:"username" example:"bagas-na@github.com"`
			}{
				Id:       resp.User.Id,
				Username: resp.User.Username,
			},
		},
	)
}

type subscriptionHandler struct {
	subscriptionclient subscriptionv1.Subscription_ServiceClient
}

func NewHandSubs(subscriptionclient subscriptionv1.Subscription_ServiceClient) *subscriptionHandler {
	return &subscriptionHandler{
		subscriptionclient: subscriptionclient,
	}
}

// GetCountries godoc
// @Summary      List available countries
// @Description  Returns countries with food security data available.
// @Tags         Countries
// @Produce      json
// @Success      200 {object} GetCountriesResponse
// @Failure      500 {object} ErrorResponse
// @Router       /countries [get]
func (h *subscriptionHandler) GetCountries(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	resp, err := h.subscriptionclient.Get_Countries(ctx, &pb.Empty{})
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	var countries []CountryInfo
	for _, c := range resp.Countries {
		country := CountryInfo{
			Name:                      c.Name,
			IpcPhase:                  c.IpcPhase,
			PopulationInPhase:         c.PopulationInPhase,
			PopulationFractionInPhase: c.PopulationFractionInPhase,
			LocationCode:              c.LocationCode,
		}

		countries = append(countries, country)
	}

	return c.JSON(http.StatusOK, GetCountriesResponse{
		Message: resp.Message,
		Data:    countries,
	})
}

// CreateSub godoc
// @Summary      Subscribe to a country
// @Description  Subscribes the authenticated user to a country.
// @Tags         Subscriptions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body CreateSubcriptionRequest true "Subscription request"
// @Success      200 {object} CreateSubcriptionResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /subscriptions [post]
func (h *subscriptionHandler) CreateSub(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	userId := c.Get("user_id").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "userId is required"})
	}

	var req CreateSubcriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	resp, err := h.subscriptionclient.Create_Subscription(ctx,
		&subscriptionv1.Subscription_Request{
			UserId:      userId,
			CountryCode: req.CountryCode,
		})
	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, CreateSubcriptionResponse{
		Message: resp.Message,
	})
}

// GetUserSubs godoc
// @Summary      Get user subscriptions
// @Description  Returns all countries subscribed by the user.
// @Tags         Subscriptions
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} GetUserSubcriptionsResponse
// @Failure      401 {object} ErrorResponse
// @Router       /subscriptions [get]
func (h *subscriptionHandler) GetUserSubs(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	userId := c.Get("user_id").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "userId is required"})
	}

	resp, err := h.subscriptionclient.Get_Subscriptions(ctx,
		&pb.Subscription_Request{
			UserId: userId,
		})

	if err != nil {
		return utils.MapGRPCError(c, err)
	}

	var subscription []SubscriptionInfo

	for _, s := range resp.Subscription {
		sub := SubscriptionInfo{
			Id:          s.Id,
			UserId:      s.UserId,
			CountryCode: s.CountryCode,
		}
		subscription = append(subscription, sub)
	}

	return c.JSON(http.StatusOK, GetUserSubcriptionsResponse{
		Data: subscription,
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

// DeleteSub godoc
// @Summary      Delete a subscription
// @Description  Unsubscribes the user from a country.
// @Tags         Subscriptions
// @Security     BearerAuth
// @Param        id path string true "Subscription ID (UUID)"
// @Success      200 {object} DeleteSubcriptionResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /subscriptions/{id} [delete]
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
