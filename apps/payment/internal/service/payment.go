package service

import (
	"context"
	"errors"
	"fmt"
	"payment-service/internal/adapters/db"
	"payment-service/internal/adapters/notification"
	stripeAdapter "payment-service/internal/adapters/stripe"

	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePaymentAndCheckout(ctx context.Context, userID uuid.UUID, countryCode string, amount int64, currency string) (*db.Payment, string, error)
	GetCheckoutURL(ctx context.Context, userID uuid.UUID, paymentID uuid.UUID) (string, error)
	ListPaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error)
	ListActivePaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error)
	UpdatePaymentToPending(ctx context.Context, paymentID uuid.UUID, providerSessionID string) (*db.Payment, error)
	UpdatePaymentToPaid(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error)
	UpdatePaymentToFailed(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error)
	UpdatePaymentToExpired(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error)
}

type paymentService struct {
	repo       *db.PaymentRepo
	stripe     *stripeAdapter.StripeAdapter
	mailer     notification.Mailer
	successURL string
	cancelURL  string
}

func NewPaymentService(repo *db.PaymentRepo, client *stripeAdapter.StripeAdapter, mailer notification.Mailer) PaymentService {
	return &paymentService{
		repo:       repo,
		stripe:     client,
		successURL: "https://example.com/success",
		cancelURL:  "https://example.com/cancel",
		mailer:     mailer,
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

func (s *paymentService) CreatePaymentAndCheckout(
	ctx context.Context,
	userID uuid.UUID,
	countryCode string,
	amount int64,
	currency string,
) (*db.Payment, string, error) {

	if amount <= 0 {
		return nil, "", errors.New("amount must be positive")
	}

	payment := &db.Payment{
		UserID:          userID,
		CountryCode:     countryCode,
		TransactionType: "payment",
		Amount:          amount,
		Currency:        currency,
		Provider:        "stripe",
		Status:          db.StatusCreated,
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

	paymentWithUser, err := s.repo.FindPaymentByID(ctx, newPayment.ID)

	fmt.Println("=== Data to be sent to notification service ===")
	fmt.Printf("payment (with user data): %#v\n", paymentWithUser)
	fmt.Printf("checkout URL: %#v\n", session.URL)

	err = s.mailer.SendCheckoutURL(ctx, paymentWithUser, session.URL)
	if err != nil {
		fmt.Printf("Error sending checkout email to %s\n", paymentWithUser.User.Username)
	}

	newPayment, err = s.UpdatePaymentToPending(ctx, newPayment.ID, session.ID)
	if err != nil {
		return nil, "", err
	}

	return newPayment, session.URL, nil
}

func (s *paymentService) GetCheckoutURL(ctx context.Context, userID, paymentID uuid.UUID) (string, error) {
	if paymentID == uuid.Nil {
		return "", errors.New("paymentID is required")
	}

	payment, err := s.repo.FindUserPaymentByID(ctx, userID, paymentID)
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

func (s *paymentService) ListPaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID is required")
	}

	return s.repo.ListByUser(ctx, userID)
}

func (s *paymentService) ListActivePaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error) {
	if userID == uuid.Nil {
		return nil, errors.New("userID is required")
	}

	return s.repo.ListPendingByUser(ctx, userID)
}

func (s *paymentService) UpdatePaymentToPending(ctx context.Context, paymentID uuid.UUID, providerSessionID string) (*db.Payment, error) {
	if providerSessionID == "" {
		return nil, errors.New("providerSessionID is required")
	}

	payment, err := s.repo.FindPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if payment.Status != db.StatusCreated {
		return nil, errors.New("payment not in created state")
	}

	return s.repo.UpdatePayment(
		ctx,
		paymentID,
		db.StatusPending,
		db.EventPending,
		providerSessionID,
		"app",
		"",
	)
}

func (s *paymentService) UpdatePaymentToPaid(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error) {
	payment, err := s.repo.FindPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if payment.Status == db.StatusPaid {
		return payment, nil
	}

	if payment.Status != db.StatusPending {
		return nil, errors.New("payment not pending")
	}

	return s.repo.UpdatePayment(
		ctx,
		paymentID,
		db.StatusPaid,
		db.EventPaid,
		"",
		"webhook",
		providerEventID,
	)
}

func (s *paymentService) UpdatePaymentToFailed(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error) {

	payment, err := s.repo.FindPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if payment.Status == db.StatusFailed {
		return payment, nil
	}

	if payment.Status != db.StatusPending {
		return nil, errors.New("payment not pending")
	}

	return s.repo.UpdatePayment(
		ctx,
		paymentID,
		db.StatusFailed,
		db.EventFailed,
		"",
		"webhook",
		providerEventID,
	)
}

func (s *paymentService) UpdatePaymentToExpired(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error) {
	payment, err := s.repo.FindPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if payment.Status == db.StatusExpired {
		return payment, nil
	}

	if payment.Status != db.StatusPending {
		return nil, errors.New("payment not pending")
	}

	return s.repo.UpdatePayment(
		ctx,
		paymentID,
		db.StatusExpired,
		db.EventExpired,
		"",
		"webhook",
		providerEventID,
	)
}
