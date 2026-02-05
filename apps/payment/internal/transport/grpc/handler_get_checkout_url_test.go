package grpcHandler_test

import (
	"context"
	paymentv1 "hunger4data/pb/payment"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetPaymentCheckoutURL_Success(t *testing.T) {
	mockSvc := new(MockPaymentService)
	paymentID := uuid.New()
	userID := uuid.New()

	mockSvc.On("GetCheckoutURL", userID, paymentID).Return("https://checkout.example.com", nil)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	resp, err := client.GetPaymentCheckoutURL(context.Background(), &paymentv1.GetPaymentCheckoutURLRequest{
		UserId:    userID.String(),
		PaymentId: paymentID.String(),
	})

	assert.NoError(t, err)
	assert.Equal(t, "https://checkout.example.com", resp.CheckoutUrl)
	mockSvc.AssertExpectations(t)
}

func TestGetPaymentCheckoutURL_InvalidPaymentID(t *testing.T) {
	mockSvc := new(MockPaymentService)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	_, err := client.GetPaymentCheckoutURL(context.Background(), &paymentv1.GetPaymentCheckoutURLRequest{
		PaymentId: "invalid-id",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockSvc.AssertNotCalled(t, "GetCheckoutURL")
}

func TestGetPaymentCheckoutURL_ServiceError(t *testing.T) {
	mockSvc := new(MockPaymentService)
	paymentID := uuid.New()
	userID := uuid.New()

	mockSvc.On("GetCheckoutURL", userID, paymentID).Return("", assert.AnError)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	_, err := client.GetPaymentCheckoutURL(context.Background(), &paymentv1.GetPaymentCheckoutURLRequest{
		PaymentId: paymentID.String(),
		UserId:    userID.String(),
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
	mockSvc.AssertExpectations(t)
}
