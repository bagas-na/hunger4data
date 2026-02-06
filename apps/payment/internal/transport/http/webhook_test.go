package httpHandler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payment-service/internal/adapters/db"
	httpHandler "payment-service/internal/transport/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v84"
)

// MockPaymentService implements db.PaymentService
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) CreatePaymentAndCheckout(ctx context.Context, userID uuid.UUID, countryCode string, amount int64, currency string) (*db.Payment, string, error) {
	args := m.Called(ctx, userID, countryCode, amount, currency)
	return args.Get(0).(*db.Payment), args.String(1), args.Error(2)
}
func (m *MockPaymentService) GetCheckoutURL(ctx context.Context, userID, paymentID uuid.UUID) (string, error) {
	args := m.Called(ctx, userID, paymentID)
	return args.String(0), args.Error(1)
}
func (m *MockPaymentService) ListPaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Payment), args.Error(1)
}
func (m *MockPaymentService) ListActivePaymentsByUser(ctx context.Context, userID uuid.UUID) ([]db.Payment, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]db.Payment), args.Error(1)
}
func (m *MockPaymentService) UpdatePaymentToPending(ctx context.Context, paymentID uuid.UUID, providerSessionID string) (*db.Payment, error) {
	args := m.Called(ctx, paymentID, providerSessionID)
	return args.Get(0).(*db.Payment), args.Error(1)
}
func (m *MockPaymentService) UpdatePaymentToPaid(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error) {
	args := m.Called(ctx, paymentID, providerEventID)
	return args.Get(0).(*db.Payment), args.Error(1)
}
func (m *MockPaymentService) UpdatePaymentToFailed(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error) {
	args := m.Called(ctx, paymentID, providerEventID)
	return args.Get(0).(*db.Payment), args.Error(1)
}
func (m *MockPaymentService) UpdatePaymentToExpired(ctx context.Context, paymentID uuid.UUID, providerEventID string) (*db.Payment, error) {
	args := m.Called(ctx, paymentID, providerEventID)
	return args.Get(0).(*db.Payment), args.Error(1)
}

func TestStripeWebhookHandler_FullPath(t *testing.T) {
	secret := "whsec_test"
	mockSvc := new(MockPaymentService)

	tests := []struct {
		name           string
		eventType      stripe.EventType
		clientRefID    string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:        "checkout.session.completed",
			eventType:   "checkout.session.completed",
			clientRefID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func() {
				mockSvc.On(
					"UpdatePaymentToPaid",
					mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					"evt_123",
				).Return(&db.Payment{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "checkout.session.expired",
			eventType:   "checkout.session.expired",
			clientRefID: "550e8400-e29b-41d4-a716-446655440001",
			mockSetup: func() {
				mockSvc.On(
					"UpdatePaymentToExpired",
					mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
					"evt_456",
				).Return(&db.Payment{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "payment_intent.payment_failed",
			eventType:   "payment_intent.payment_failed",
			clientRefID: "550e8400-e29b-41d4-a716-446655440002",
			mockSetup: func() {
				mockSvc.On(
					"UpdatePaymentToFailed",
					mock.Anything,
					uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
					"evt_789",
				).Return(&db.Payment{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Compute event ID based on test case to match mock
			eventID := "evt_" + tt.clientRefID[len(tt.clientRefID)-3:]

			// Setup mock service calls with correct eventID
			switch tt.eventType {
			case "checkout.session.completed":
				mockSvc.On(
					"UpdatePaymentToPaid",
					mock.Anything,
					uuid.MustParse(tt.clientRefID),
					eventID,
				).Return(&db.Payment{}, nil)
			case "checkout.session.expired":
				mockSvc.On(
					"UpdatePaymentToExpired",
					mock.Anything,
					uuid.MustParse(tt.clientRefID),
					eventID,
				).Return(&db.Payment{}, nil)
			case "payment_intent.payment_failed":
				mockSvc.On(
					"UpdatePaymentToFailed",
					mock.Anything,
					uuid.MustParse(tt.clientRefID),
					eventID,
				).Return(&db.Payment{}, nil)
			}

			// Override ConstructEvent to bypass real Stripe signature
			httpHandler.ConstructEvent = func(payload []byte, sigHeader, secret string, opts ...stripe.WebhookOption) (stripe.Event, error) {
				var raw []byte
				if tt.eventType == "payment_intent.payment_failed" {
					obj := map[string]map[string]string{
						"metadata": {"payment_id": tt.clientRefID},
					}
					raw, _ = json.Marshal(obj)
				} else {
					session := stripe.CheckoutSession{
						ClientReferenceID: tt.clientRefID,
					}
					raw, _ = json.Marshal(session)
				}
				return stripe.Event{
					ID:   eventID, // must match mock
					Type: tt.eventType,
					Data: &stripe.EventData{Raw: raw},
				}, nil
			}
			defer func() { httpHandler.ConstructEvent = stripe.ConstructEvent }() // restore

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`dummy payload`)))
			req.Header.Set("Stripe-Signature", "ignored")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := httpHandler.StripeWebhookHandler(secret, mockSvc)
			err := handler(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockSvc.ExpectedCalls = nil // reset mocks
		})
	}

}
