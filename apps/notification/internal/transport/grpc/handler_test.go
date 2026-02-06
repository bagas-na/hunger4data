package grpcHandler

import (
	"context"
	"errors"
	"testing"

	notifyv1 "hunger4data/pb/notification"
	"notification/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// If your EmailService isn't an interface, you can mock it by wrapping it
// or using a library like testify/mock.
// For this example, I'll assume a mockable structure.

type MockEmailService struct {
	mock.Mock
}

// We simulate the SendNotification method
func (m *MockEmailService) SendNotification(ctx context.Context, email service.Email) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func TestSendTransactionEmail(t *testing.T) {
	fromName := "Test System"
	fromEmail := "system@test.com"

	tests := []struct {
		name          string
		request       *notifyv1.SendEmailRequest
		mockBehavior  func(m *MockEmailService)
		expectedResp  *notifyv1.SendEmailResponse
		expectedError error
	}{
		{
			name: "Success - Email Sent",
			request: &notifyv1.SendEmailRequest{
				To:      "user@example.com",
				Subject: "Welcome",
				Body:    "Hello there!",
			},
			mockBehavior: func(m *MockEmailService) {
				m.On("SendNotification", mock.Anything, mock.MatchedBy(func(e service.Email) bool {
					return e.ToEmail == "user@example.com" && e.Subject == "Welcome"
				})).Return(nil)
			},
			expectedResp: &notifyv1.SendEmailResponse{Success: true},
		},
		{
			name: "Failure - Service Error",
			request: &notifyv1.SendEmailRequest{
				To: "fail@example.com",
			},
			mockBehavior: func(m *MockEmailService) {
				m.On("SendNotification", mock.Anything, mock.Anything).
					Return(errors.New("provider down"))
			},
			expectedResp: &notifyv1.SendEmailResponse{
				Success: false,
				Error:   "provider down",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := new(MockEmailService)
			tt.mockBehavior(mockSvc)

			// Note: If your handler expects the concrete *service.EmailService,
			// you'll need to use an interface in your handler struct instead.
			s := &EmailGRPCServer{
				svc:       mockSvc, // Ensure handler.go uses an interface for this to work!
				fromName:  fromName,
				fromEmail: fromEmail,
			}

			// Execute
			resp, err := s.SendTransactionEmail(context.Background(), tt.request)

			// Assert
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResp, resp)
			mockSvc.AssertExpectations(t)
		})
	}
}
