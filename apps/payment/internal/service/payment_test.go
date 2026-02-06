package service_test

import (
	"context"
	"payment-service/internal/adapters/db"
	stripeAdapter "payment-service/internal/adapters/stripe"
	"payment-service/internal/mocks"
	"payment-service/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Tests

func TestCreatePaymentAndCheckout(t *testing.T) {
	mockRepo := new(mocks.PaymentRepo)
	mockStripe := new(mocks.StripeAdapter)
	mockMailer := new(mocks.Mailer)
	svc := service.NewPaymentService(mockRepo, mockStripe, mockMailer)

	userID := uuid.New()
	paymentID := uuid.New()

	tests := []struct {
		name        string
		countryCode string
		amount      int64
		mockSetup   func()
		wantErr     bool
	}{
		{
			name:        "invalid country code",
			countryCode: "XXX",
			amount:      1000,
			mockSetup:   func() {},
			wantErr:     true,
		},
		{
			name:        "successful payment",
			countryCode: "IDN",
			amount:      5000,
			mockSetup: func() {
				mockRepo.On("CreatePayment", mock.Anything, mock.Anything).Return(&db.Payment{
					ID:     paymentID,
					Status: db.StatusCreated,
				}, nil)

				mockStripe.On("CreateCheckoutSession", mock.Anything, mock.Anything).Return(&stripeAdapter.CheckoutSessionResult{
					ID:  "sess_123",
					URL: "https://checkout.url",
				}, nil)

				mockRepo.On("FindPaymentByID", mock.Anything, paymentID).Return(&db.Payment{
					ID:     paymentID,
					Status: db.StatusCreated, // <-- must be CREATED
					User: db.User{
						Username: "john",
					},
				}, nil)

				mockMailer.On("SendCheckoutURL", mock.Anything, mock.Anything, "https://checkout.url").Return(nil)

				mockRepo.On("UpdatePayment", mock.Anything, paymentID, db.StatusPending, db.EventPending, "sess_123", "app", "").Return(&db.Payment{
					ID:     paymentID,
					Status: db.StatusPending,
				}, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			p, url, err := svc.CreatePaymentAndCheckout(context.Background(), userID, tt.countryCode, tt.amount, "idr")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, p)
				assert.Equal(t, "https://checkout.url", url)
			}
			mockRepo.ExpectedCalls = nil
			mockStripe.ExpectedCalls = nil
			mockMailer.ExpectedCalls = nil
		})
	}
}
