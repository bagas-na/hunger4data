package httpHandler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v84"
)

func StripeWebhookHandler(webhookSecret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		event := stripe.Event{}

		if err := c.Bind(&event); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "invalid argument",
				"detail":  err.Error(),
			})
		}

		sig := c.Request().Header.Get("Stripe-Signature")
		if sig == "" {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "empty stripe signature",
			})
		}

		switch event.Type {
		case "checkout.session.async_payment_failed":
			{
				fmt.Println("checkout.session.async_payment_failed")
			}
		case "checkout.session.async_payment_succeeded":
			{
				fmt.Println("checkout.session.async_payment_succeeded")
			}
		case "checkout.session.completed":
			{
				fmt.Println("checkout.session.completed")
			}
		case "checkout.session.expired":
			{
				fmt.Println("checkout.session.expired")
			}
		case "payment_intent.canceled":
			{
				fmt.Println("payment_intent.canceled")
			}
		case "payment_intent.created":
			{
				fmt.Println("payment_intent.created")
			}
		case "payment_intent.payment_failed":
			{
				fmt.Println("payment_intent.payment_failed")
			}
		case "payment_intent.processing":
			{
				fmt.Println("payment_intent.processing")
			}
		case "payment_intent.succeeded":
			{
				fmt.Println("payment_intent.succeeded")
			}

		}
		return c.NoContent(http.StatusOK)
	}
}
