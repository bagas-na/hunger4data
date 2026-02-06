package notification

import (
	"context"
	"testing"

	notifyv1 "hunger4data/pb/notification"
	"payment-service/internal/adapters/db"
	"payment-service/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockEmailServiceClient
type MockEmailServiceClient struct {
	mock.Mock
}

func (m *MockEmailServiceClient) SendTransactionEmail(ctx context.Context, in *notifyv1.SendEmailRequest, opts ...grpc.CallOption) (*notifyv1.SendEmailResponse, error) {
	args := m.Called(ctx, in)

	var res *notifyv1.SendEmailResponse
	if args.Get(0) != nil {
		res = args.Get(0).(*notifyv1.SendEmailResponse)
	}
	return res, args.Error(1)
}

func TestSendCheckoutURL(t *testing.T) {
	utils.ISO3166Alpha3 = map[string]string{
		"US": "United States",
	}

	mockClient := new(MockEmailServiceClient)
	m := NewMailer(mockClient)

	payment := &db.Payment{
		CountryCode: "US",
		User: db.User{
			Username: "alice@example.com",
		},
	}
	checkoutURL := "https://stripe.com/checkout/123"

	mockClient.On("SendTransactionEmail", mock.Anything, mock.MatchedBy(func(req *notifyv1.SendEmailRequest) bool {
		return req.To == "alice@example.com" && assert.Contains(t, req.Body, "United States")
	})).Return(nil, nil)

	// Execute
	err := m.SendCheckoutURL(context.Background(), payment, checkoutURL)

	// Assert
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}
