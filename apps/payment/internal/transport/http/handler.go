package httpHandler

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CreatePaymentRequest struct {
	CountryID uuid.UUID `json:"country_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
}

func (h *Handler) CreatePayment(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	var req CreatePaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userID := uuid.New() // TODO: diganti dengan userID dari JWT

	payment, err := h.service.CreatePayment(
		ctx,
		userID,
		req.CountryID,
		req.Amount,
		req.Currency,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, payment)
}

func (h *Handler) GetPayment(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
	defer cancel()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	Payment, err := h.service.GetPayment(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, Payment)
}
