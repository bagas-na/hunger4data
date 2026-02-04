package grpcHandler

import (
	paymentv1 "hunger4data/pb/payment"
	"payment-service/internal/adapters/db"
	"time"
)

func mapPaymentToProto(p *db.Payment) *paymentv1.Payment {
	return &paymentv1.Payment{
		Id:                p.ID.String(),
		UserId:            p.UserID.String(),
		CountryId:         p.CountryID.String(),
		TransactionType:   p.TransactionType,
		Amount:            p.Amount,
		Currency:          p.Currency,
		Provider:          p.Provider,
		ProviderSessionId: p.ProviderSessionID,
		Status:            string(p.Status),
		CreatedAt:         p.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:         p.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
