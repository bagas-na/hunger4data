package grpcHandler

import (
	"context"
	paymentv1 "hunger4data/pb/payment"
	"payment-service/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentGRPCServer struct {
	paymentv1.UnimplementedPaymentServiceServer
	svc *service.PaymentService
}

func NewPaymentGRPCServer(svc *service.PaymentService) *PaymentGRPCServer {
	return &PaymentGRPCServer{
		svc: svc,
	}
}

func (s *PaymentGRPCServer) CreatePayment(ctx context.Context, req *paymentv1.CreatePaymentRequest) (*paymentv1.CreatePaymentResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	countryID, err := uuid.Parse(req.CountryId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid country_id")
	}

	payment, checkoutURL, err := s.svc.CreatePaymentAndCheckout(
		ctx,
		userID,
		countryID,
		req.Amount,
		req.Currency,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &paymentv1.CreatePaymentResponse{
		Payment:     mapPayment(payment),
		CheckoutUrl: checkoutURL,
	}, nil
}

func (s *PaymentGRPCServer) GetPaymentCheckoutURL(ctx context.Context, req *paymentv1.GetPaymentCheckoutURLRequest) (*paymentv1.GetPaymentCheckoutURLResponse, error) {
	paymentID, err := uuid.Parse(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment_id")
	}

	url, err := s.svc.GetCheckoutURL(ctx, paymentID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &paymentv1.GetPaymentCheckoutURLResponse{
		CheckoutUrl: url,
	}, nil
}

func (s *PaymentGRPCServer) ListPayments(ctx context.Context, req *paymentv1.ListPaymentsRequest) (*paymentv1.ListPaymentsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	payments, err := s.svc.ListPaymentsByUser(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &paymentv1.ListPaymentsResponse{}
	for _, p := range payments {
		resp.Payments = append(resp.Payments, mapPayment(&p))
	}

	return resp, nil
}

func (s *PaymentGRPCServer) ListPendingPayments(ctx context.Context, req *paymentv1.ListPendingPaymentsRequest) (*paymentv1.ListPendingPaymentsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	payments, err := s.svc.ListActivePaymentsByUser(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &paymentv1.ListPendingPaymentsResponse{}
	for _, p := range payments {
		resp.Payments = append(resp.Payments, mapPayment(&p))
	}

	return resp, nil
}
