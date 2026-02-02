package grpcHandler

import (
	paymentv1 "hunger4data/pb/payment"
	"payment-service/internal/adapters/db"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapPayment(p *db.Payment) *paymentv1.Payment {
	return &paymentv1.Payment{
		Id:                p.ID.String(),
		UserId:            p.UserID.String(),
		CountryId:         p.CountryID.String(),
		TransactionType:   p.TransactionType,
		Amount:            p.Amount,
		Currency:          p.Currency,
		Provider:          p.Provider,
		ProviderSessionId: p.ProviderSessionID,
		Status:            p.Status,
		CreatedAt:         timestamppb.New(p.CreatedAt),
		UpdatedAt:         timestamppb.New(p.UpdatedAt),
	}
}
