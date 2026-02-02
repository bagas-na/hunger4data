package service

import (
	"context"
	"errors"
	"payment-service/internal/adapters/db"
	stripeAdapter "payment-service/internal/adapters/stripe"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo       *db.PaymentRepo
	stripe     *stripeAdapter.StripeAdapter
	successURL string
	cancelURL  string
}

func NewPaymentService(repo *db.PaymentRepo, client *stripeAdapter.StripeAdapter) *PaymentService {
	return &PaymentService{
		repo:       repo,
		stripe:     client,
		successURL: "https://example.com/success",
		cancelURL:  "https://example.com/cancel",
	}
}

// func (s *PaymentService) CreatePaymentDBOnly(
// 	ctx context.Context,
// 	userID uuid.UUID,
// 	countryID uuid.UUID,
// 	amount int64,
// 	currency string,
// ) (*db.Payment, error) {

// 	if amount <= 0 {
// 		return nil, errors.New("amount must be positive")
// 	}

// 	payment := &db.Payment{
// 		UserID:          userID,
// 		CountryID:       countryID,
// 		TransactionType: "payment",
// 		Amount:          amount,
// 		Currency:        currency,
// 		Provider:        "internal",
// 	}

// 	return s.repo.CreatePayment(ctx, payment)
// }

func (s *PaymentService) CreatePaymentAndCheckout(
	ctx context.Context,
	userID uuid.UUID,
	countryID uuid.UUID,
	amount int64,
	currency string,
) (*db.Payment, string, error) {

	if amount <= 0 {
		return nil, "", errors.New("amount must be positive")
	}

	payment := &db.Payment{
		UserID:          userID,
		CountryID:       countryID,
		TransactionType: "payment",
		Amount:          amount,
		Currency:        currency,
		Provider:        "stripe",
		Status:          "pending",
	}

	newPayment, err := s.repo.CreatePayment(ctx, payment)
	if err != nil {
		return nil, "", err
	}

	session, err := s.stripe.CreateCheckoutSession(
		ctx,
		stripeAdapter.CreateCheckoutSessionInput{
			Amount:     amount,
			Currency:   currency,
			PaymentID:  newPayment.ID.String(),
			SuccessURL: s.successURL,
			CancelURL:  s.cancelURL,
		},
	)
	if err != nil {
		return nil, "", err
	}

	newPayment, err = s.repo.UpdatePaymentProviderObject(
		ctx,
		newPayment.ID,
		session.PaymentIntentID,
	)
	if err != nil {
		return nil, "", err
	}

	return newPayment, session.URL, nil
}

func (s *PaymentService) GetCheckoutURL(ctx context.Context, paymentID uuid.UUID) (string, error) {
	if paymentID == uuid.Nil {
		return "", errors.New("paymentID is required")
	}

	payment, err := s.repo.FindByID(ctx, paymentID)
	if err != nil {
		return "", err
	}

	if payment.Provider != "stripe" {
		return "", errors.New("payment is not handled by stripe")
	}

	if payment.ProviderSessionID == "" {
		return "", errors.New("checkout session not created for this payment")
	}

	return s.stripe.GetCheckoutSessionURL(
		ctx,
		payment.ProviderSessionID,
	)
}

func (s *PaymentService) ListPaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID is required")
	}

	return s.repo.ListByUser(ctx, userID)
}

func (s *PaymentService) ListActivePaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID is required")
	}

	return s.repo.ListPendingByUser(ctx, userID)
}
