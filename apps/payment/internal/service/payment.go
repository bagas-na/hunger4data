package service

import (
	"context"
	"errors"
	"payment-service/internal/adapters/db"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo *db.PaymentRepo
}

func NewPaymentService(repo *db.PaymentRepo) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) CreatePayment(
	ctx context.Context,
	userID uuid.UUID,
	countryID uuid.UUID,
	amount int64,
	currency string,
) (*db.Payment, error) {

	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	payment := &db.Payment{
		UserID:          userID,
		CountryID:       countryID,
		TransactionType: "payment",
		Amount:          amount,
		Currency:        currency,
		Provider:        "internal",
	}

	return s.repo.CreatePayment(ctx, payment)
}

func (s *PaymentService) GetPayment(
	ctx context.Context,
	id uuid.UUID,
) (*db.Payment, error) {

	return s.repo.FindByID(ctx, id)
}
