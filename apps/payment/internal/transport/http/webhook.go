package httpHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payment-service/internal/service"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v84"
)

func StripeWebhookHandler(webhookSecret string, svc service.PaymentService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
		defer cancel()

		payload, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "failed to read request body",
				"detail":  err.Error(),
			})
		}

		sig := c.Request().Header.Get("Stripe-Signature")
		if sig == "" {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "missing Stripe-Signature header",
			})
		}

		event, err := stripe.ConstructEvent(payload, sig, webhookSecret)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "invalid stripe signature",
			})
		}

		switch event.Type {

		case "checkout.session.completed", "checkout.session.expired":
			{
				var session stripe.CheckoutSession
				if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
					return c.JSON(http.StatusBadRequest, map[string]any{
						"message": "unable to unmarshal CheckoutSession",
					})
				}

				paymentID, err := uuid.Parse(session.ClientReferenceID)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]any{
						"message": "invalid ClientReferenceID",
					})
				}

				if event.Type == "checkout.session.completed" {
					_, err = svc.UpdatePaymentToPaid(ctx, paymentID, event.ID)
					fmt.Printf("\nUPDATING PAYMENT %q TO [PAID] (event: %q)\n\n", paymentID, event.ID)
				} else {
					_, err = svc.UpdatePaymentToExpired(ctx, paymentID, event.ID)
					fmt.Printf("\nUPDATING PAYMENT %q TO [EXPIRED] (event: %q)\n\n", paymentID, event.ID)
				}

				if err != nil {
					return err
				}
			}

		case "payment_intent.payment_failed":
			{
				var pi stripe.PaymentIntent
				if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
					return c.JSON(http.StatusBadRequest, map[string]any{
						"message": "unable to unmarshal PaymentIntent",
					})
				}

				paymentID, err := uuid.Parse(pi.Metadata["payment_id"])
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]any{
						"message": "invalid payment_id metadata",
					})
				}

				if _, err := svc.UpdatePaymentToFailed(ctx, paymentID, event.ID); err != nil {
					return err
				}
				fmt.Printf("\nUPDATING PAYMENT %q TO [FAILED] (event: %q)\n\n", paymentID, event.ID)
			}
		}

		return c.NoContent(http.StatusOK)
	}
}
