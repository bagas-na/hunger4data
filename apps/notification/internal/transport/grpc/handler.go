package grpcHandler

import (
	"context"
	notifyv1 "hunger4data/pb/notification"
	"notification/internal/service"
)

// type EmailGRPCServer interface {
// 	SendTransactionEmail(ctx context.Context, req *notifyv1.SendEmailRequest) (*notifyv1.SendEmailResponse, error)
// }

type EmailGRPCServer struct {
	notifyv1.UnimplementedEmailServiceServer
	svc       service.EmailService
	fromName  string
	fromEmail string
}

func NewEmailGRPCServer(svc service.EmailService, fromName, fromEmail string) *EmailGRPCServer {
	return &EmailGRPCServer{
		svc:       svc,
		fromName:  fromName,
		fromEmail: fromEmail,
	}
}

func (s *EmailGRPCServer) SendTransactionEmail(ctx context.Context, req *notifyv1.SendEmailRequest) (*notifyv1.SendEmailResponse, error) {
	email := service.Email{
		FromName:  s.fromName,
		FromEmail: s.fromEmail,
		ToName:    req.To,
		ToEmail:   req.To,
		Subject:   req.Subject,
		Body:      req.Body,
	}

	if err := s.svc.SendNotification(ctx, email); err != nil {
		return &notifyv1.SendEmailResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &notifyv1.SendEmailResponse{
		Success: true,
	}, nil
}
