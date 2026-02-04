package grpcHandler_test

import (
	"context"
	paymentv1 "hunger4data/pb/payment"
	"payment-service/internal/adapters/db"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Tests
func TestCreatePayment_Success(t *testing.T) {
	mockSvc := new(MockPaymentService)

	userID := uuid.New()
	countryCode := "AFG"
	paymentID := uuid.New()

	mockSvc.On(
		"CreatePaymentAndCheckout",
		userID,
		countryCode,
		int64(1000),
		"USD",
	).Return(
		&db.Payment{
			ID:          paymentID,
			UserID:      userID,
			CountryCode: countryCode,
			Amount:      1000,
			Currency:    "USD",
			Status:      "pending",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		"https://checkout.stripe.com/test",
		nil,
	)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	resp, err := client.CreatePayment(context.Background(), &paymentv1.CreatePaymentRequest{
		UserId:      userID.String(),
		CountryCode: countryCode,
		Amount:      1000,
		Currency:    "USD",
	})

	assert.NoError(t, err)
	assert.Equal(t, "pending", resp.Payment.Status)
	assert.NotEmpty(t, resp.CheckoutUrl)

	mockSvc.AssertExpectations(t)
}

func TestCreatePayment_InvalidUserID(t *testing.T) {
	mockSvc := new(MockPaymentService)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	_, err := client.CreatePayment(context.Background(), &paymentv1.CreatePaymentRequest{
		UserId:      "not-a-uuid",
		CountryCode: uuid.New().String(),
		Amount:      1000,
		Currency:    "USD",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	mockSvc.AssertNotCalled(t, "CreatePaymentAndCheckout")
}

func TestCreatePayment_ServiceError(t *testing.T) {
	mockSvc := new(MockPaymentService)

	userID := uuid.New()
	countryID := uuid.New()

	mockSvc.On(
		"CreatePaymentAndCheckout",
		userID,
		countryID,
		int64(1000),
		"USD",
	).Return(
		(*db.Payment)(nil),
		"",
		assert.AnError,
	)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	_, err := client.CreatePayment(context.Background(), &paymentv1.CreatePaymentRequest{
		UserId:      userID.String(),
		CountryCode: countryID.String(),
		Amount:      1000,
		Currency:    "USD",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
}
