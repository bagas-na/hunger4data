package grpcHandler_test

import (
	"context"
	paymentv1 "hunger4data/pb/payment"
	"payment-service/internal/adapters/db"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestListPendingPayments_Success(t *testing.T) {
	mockSvc := new(MockPaymentService)
	userID := uuid.New()
	payments := []db.Payment{
		{ID: uuid.New(), Amount: 100},
		{ID: uuid.New(), Amount: 200},
	}

	mockSvc.On("ListActivePaymentsByUser", userID).Return(payments, nil)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	resp, err := client.ListPendingPayments(context.Background(), &paymentv1.ListPendingPaymentsRequest{
		UserId: userID.String(),
	})

	assert.NoError(t, err)
	assert.Len(t, resp.Payments, 2)
	mockSvc.AssertExpectations(t)
}

func TestListPendingPayments_InvalidUserID(t *testing.T) {
	mockSvc := new(MockPaymentService)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	_, err := client.ListPendingPayments(context.Background(), &paymentv1.ListPendingPaymentsRequest{
		UserId: "invalid-id",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockSvc.AssertNotCalled(t, "ListActivePaymentsByUser")
}

func TestListPendingPayments_ServiceError(t *testing.T) {
	mockSvc := new(MockPaymentService)
	userID := uuid.New()

	mockSvc.On("ListActivePaymentsByUser", userID).
		Return([]db.Payment(nil), assert.AnError)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	_, err := client.ListPendingPayments(context.Background(), &paymentv1.ListPendingPaymentsRequest{
		UserId: userID.String(),
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	mockSvc.AssertExpectations(t)
}
