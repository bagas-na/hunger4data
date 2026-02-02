package stripeAdapter

import (
	"context"

	"github.com/stripe/stripe-go/v84"
)

type StripeAdapter struct {
	client *stripe.Client
}

func NewStripeAdapter(secretKey string) *StripeAdapter {
	return &StripeAdapter{
		client: stripe.NewClient(secretKey),
	}
}

// type PaymentIntentResult struct {
// 	ID           string
// 	ClientSecret string
// 	Status       string
// }

// func (a *StripeAdapter) CreatePaymentIntent(
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
	ID  string
	URL string
}

func (a *StripeAdapter) CreateCheckoutSession(
	ctx context.Context,
	input CreateCheckoutSessionInput,
) (*CheckoutSessionResult, error) {

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
	}

	s, err := a.client.V1CheckoutSessions.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	return &CheckoutSessionResult{
		ID:  s.ID,
		URL: s.URL,
	}, nil
}
