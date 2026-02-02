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

func TestListPayments_Success(t *testing.T) {
	mockSvc := new(MockPaymentService)
	userID := uuid.New()
	payments := []db.Payment{
		{ID: uuid.New(), Amount: 100},
		{ID: uuid.New(), Amount: 200},
	}

	mockSvc.On("ListPaymentsByUser", userID).Return(payments, nil)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	resp, err := client.ListPayments(context.Background(), &paymentv1.ListPaymentsRequest{
		UserId: userID.String(),
	})

	assert.NoError(t, err)
	assert.Len(t, resp.Payments, 2)
	mockSvc.AssertExpectations(t)
}

func TestListPayments_InvalidUserID(t *testing.T) {
	mockSvc := new(MockPaymentService)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	_, err := client.ListPayments(context.Background(), &paymentv1.ListPaymentsRequest{
		UserId: "invalid-id",
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockSvc.AssertNotCalled(t, "ListPaymentsByUser")
}

func TestListPayments_ServiceError(t *testing.T) {
	mockSvc := new(MockPaymentService)
	userID := uuid.New()

	mockSvc.On("ListPaymentsByUser", userID).
		Return([]db.Payment(nil), assert.AnError)

	conn, cleanup := startTestGRPCServer(t, mockSvc)
	defer cleanup()

	client := paymentv1.NewPaymentServiceClient(conn)

	_, err := client.ListPayments(context.Background(), &paymentv1.ListPaymentsRequest{
		UserId: userID.String(),
	})

	st, _ := status.FromError(err)
	assert.Equal(t, codes.Internal, st.Code())
	mockSvc.AssertExpectations(t)
}
