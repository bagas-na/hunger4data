package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{
		db: db,
	}
}

func (r *PaymentRepo) FindByID(ctx context.Context, id uuid.UUID) (*Payment, error) {
	var payment Payment
	err := r.db.WithContext(ctx).First(&payment, "id = ?", id).Error
	return &payment, err
}

func (r *PaymentRepo) CreatePayment(ctx context.Context, p *Payment) (*Payment, error) {
	paymentID := uuid.New()
	eventID := uuid.New()

	newPayment := *p
	newPayment.ID = paymentID
	newPayment.Status = "pending"
	newPayment.CreatedAt = time.Now()
	newPayment.UpdatedAt = time.Now()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newPayment).Error; err != nil {
			return err
		}

		newEvent := PaymentEvent{
			ID:        eventID,
			PaymentID: newPayment.ID,
			EventType: "created",
			Source:    "app",
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&newEvent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &newPayment, nil
}

func (r *PaymentRepo) UpdatePayment(
	ctx context.Context,
	id uuid.UUID,
	status string,
	eventType string,
	objectStr string,
	source string,
) (*Payment, error) {

	eventID := uuid.New()
	var payment Payment

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.First(&payment, "id = ?", id).Error; err != nil {
			return err
		}

		if status != "" {
			payment.Status = status
		}

		if objectStr != "" {
			payment.ProviderObjectID = objectStr
		}

		payment.UpdatedAt = time.Now()

		if err := tx.Save(&payment).Error; err != nil {
			return err
		}

		newEvent := PaymentEvent{
			ID:        eventID,
			PaymentID: payment.ID,
			EventType: eventType,
			Source:    source,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&newEvent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &payment, nil
}
