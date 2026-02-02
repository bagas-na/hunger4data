package grpcHandler_test

import (
	"context"
	"net"
	"testing"

	"payment-service/internal/adapters/db"
	"payment-service/internal/service"
	grpcHandler "payment-service/internal/transport/grpc"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentv1 "hunger4data/pb/payment"
)

// Mocks Interface
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) CreatePaymentAndCheckout(
	ctx context.Context,
	userID uuid.UUID,
	countryID uuid.UUID,
	amount int64,
	currency string,
) (*db.Payment, string, error) {

	args := m.Called(userID, countryID, amount, currency)
	return args.Get(0).(*db.Payment), args.String(1), args.Error(2)
}

func (m *MockPaymentService) GetCheckoutURL(ctx context.Context, paymentID uuid.UUID) (string, error) {
	args := m.Called(paymentID)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentService) ListPaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error) {
	args := m.Called(userID)
	return args.Get(0).([]db.Payment), args.Error(1)
}

func (m *MockPaymentService) ListActivePaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error) {
	args := m.Called(userID)
	return args.Get(0).([]db.Payment), args.Error(1)
}

func (m *MockPaymentService) UpdatePaymentToPending(ctx context.Context, paymentID uuid.UUID, providerSessionID string) (*db.Payment, error) {
	args := m.Called(paymentID, providerSessionID)
	return args.Get(0).(*db.Payment), args.Error(1)
}

func (m *MockPaymentService) UpdatePaymentToPaid(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error) {
	args := m.Called(paymentID, providerEventID)
	return args.Get(0).(*db.Payment), args.Error(1)
}

func (m *MockPaymentService) UpdatePaymentToFailed(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error) {
	args := m.Called(paymentID, providerEventID)
	return args.Get(0).(*db.Payment), args.Error(1)
}

func (m *MockPaymentService) UpdatePaymentToExpired(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error) {
	args := m.Called(paymentID, providerEventID)
	return args.Get(0).(*db.Payment), args.Error(1)
}

// Helper
func startTestGRPCServer(t *testing.T, svc service.PaymentService) (*grpc.ClientConn, func()) {
	t.Helper()

	lis, err := net.Listen("tcp", ":56565")
	if err != nil {
		t.Fatal(err)
	}

	server := grpc.NewServer()
	handler := grpcHandler.NewPaymentGRPCServer(svc)
	paymentv1.RegisterPaymentServiceServer(server, handler)

	go server.Serve(lis)

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		conn.Close()
		server.Stop()
	}

	return conn, cleanup
}
