package grpc

import (
	"authenticator/internal/service"
	proto "authenticator/proto/authenticator"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authclient proto.AuthServiceClient
}

func NewHandAuth(authclient proto.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		authclient: authclient,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req service.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	ctx := context.TODO()

	resp, err := h.authclient.Login(ctx, &proto.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "error logging in",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": resp.Token,
	})
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req service.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	resp, err := h.authclient.Register(ctx, &proto.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "error registering",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": resp.Message,
		"user": map[string]interface{}{
			"id":       resp.User.Id,
			"username": resp.User.Username,
		},
	})
}
