package stripeAdapter

import (
	"context"
	"errors"

	"github.com/stripe/stripe-go/v84"
)

type StripeAdapter interface {
	CreateCheckoutSession(ctx context.Context, input CreateCheckoutSessionInput) (*CheckoutSessionResult, error)
	GetCheckoutSessionURL(ctx context.Context, sessionID string) (string, error)
}

type stripeAdapter struct {
	client *stripe.Client
}

func NewStripeAdapterWithClient(client *stripe.Client) StripeAdapter {
	return &stripeAdapter{client: client}
}

// type PaymentIntentResult struct {
// 	ID           string
// 	ClientSecret string
// 	Status       string
// }

// func (a *stripeAdapter) CreatePaymentIntent(
// 	ctx context.Context,
// 	amount int64,
// 	currency string,
// 	idempotencyKey string,
// 	metadata map[string]string,
// ) (*PaymentIntentResult, error) {

// 	params := &stripe.PaymentIntentCreateParams{
// 		Amount:   stripe.Int64(amount),
// 		Currency: stripe.String(currency),
// 	}

// 	pi, err := a.client.V1PaymentIntents.Create(ctx, params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &PaymentIntentResult{
// 		ID:           pi.ID,
// 		ClientSecret: pi.ClientSecret,
// 		Status:       string(pi.Status),
// 	}, nil
// }

type CreateCheckoutSessionInput struct {
	Amount     int64
	Currency   string
	PaymentID  string
	SuccessURL string
	CancelURL  string
}

type CheckoutSessionResult struct {
	ID              string
	PaymentIntentID string
	URL             string
}

func (a *stripeAdapter) CreateCheckoutSession(ctx context.Context, input CreateCheckoutSessionInput) (*CheckoutSessionResult, error) {
	params := &stripe.CheckoutSessionCreateParams{
		Mode:       stripe.String("payment"),
		SuccessURL: stripe.String(input.SuccessURL),
		CancelURL:  stripe.String(input.CancelURL),

		ClientReferenceID: stripe.String(input.PaymentID),
		Metadata: map[string]string{
			"payment_id": input.PaymentID,
		},

		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
					Currency:   stripe.String(input.Currency),
					UnitAmount: stripe.Int64(input.Amount),
					ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
						Name: stripe.String("Donation"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},

		Expand: []*string{
			stripe.String("payment_intent"),
		},
	}

	s, err := a.client.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	var piID string
	if s.PaymentIntent != nil {
		piID = s.PaymentIntent.ID
	}

	return &CheckoutSessionResult{
		ID:              s.ID,
		PaymentIntentID: piID,
		URL:             s.URL,
	}, nil
}

func (a *stripeAdapter) GetCheckoutSessionURL(ctx context.Context, sessionID string) (string, error) {
	s, err := a.client.V1CheckoutSessions.Retrieve(ctx, sessionID, nil)
	if err != nil {
		return "", err
	}

	if s.URL == "" {
		return "", errors.New("checkout session has no URL")
	}

	return s.URL, nil
}
