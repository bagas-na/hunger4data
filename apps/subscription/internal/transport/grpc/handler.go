package grpc

import (
	"context"
	"net/http"
	"subscription/internal/adapters/model"
	proto "subscription/proto/subcription"

	"github.com/labstack/echo"
)

type subscriptionHandler struct {
	subscriptionclient proto.Subscription_ServiceClient
}

func NewHandAuth(subscriptionclient proto.Subscription_ServiceClient) *subscriptionHandler {
	return &subscriptionHandler{
		subscriptionclient: subscriptionclient,
	}
}

func (h *subscriptionHandler) GetCountries(c echo.Context) error {
	var req model.Country
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}
	ctx := context.TODO()
	resp, err := h.subscriptionclient.Get_Countries(ctx, &proto.Empty{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": resp.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": resp.Countries,
	})
}

func (h *subscriptionHandler) CreateSub(c echo.Context) error {
	var req model.Subscription
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}
	ctx := context.TODO()
	resp, err := h.subscriptionclient.Create_Subscription(ctx, &proto.Subscription_Request{
		IdUser:    req.Id_user,
		IdCountry: req.Id_country,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": resp.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": resp.Message,
	})
}

func (h *subscriptionHandler) GetSubByID(c echo.Context) error {
	var req model.Subscription
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}
	ctx := context.TODO()
	resp, err := h.subscriptionclient.Get_Subscription_By_ID(ctx, &proto.Subscription_Request{
		Id:        req.Id,
		IdUser:    req.Id_user,
		IdCountry: req.Id_country,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": resp.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": resp.Message,
	})
}

func (h *subscriptionHandler) UpdateSub(c echo.Context) error {
	var req model.Subscription
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}
	ctx := context.TODO()
	resp, err := h.subscriptionclient.Update_Subscription(ctx, &proto.Subscription_Request{
		Id:        req.Id,
		IdUser:    req.Id_user,
		IdCountry: req.Id_country,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": resp.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": resp.Message,
	})
}

func (h *subscriptionHandler) DeleteSub(c echo.Context) error {
	var req model.Subscription
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}
	ctx := context.TODO()
	resp, err := h.subscriptionclient.Delete_Subscription(ctx, &proto.Subscription_Request{
		Id:        req.Id,
		IdUser:    req.Id_user,
		IdCountry: req.Id_country,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": resp.Message,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": resp.Message,
	})
}
