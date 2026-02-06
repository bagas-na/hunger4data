package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepo interface {
	CreatePayment(ctx context.Context, p *Payment) (*Payment, error)
	FindPaymentByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	FindUserPaymentByID(ctx context.Context, userId uuid.UUID, paymentId uuid.UUID) (*Payment, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Payment, error)
	ListPendingByUser(ctx context.Context, userID uuid.UUID) ([]Payment, error)
	UpdatePayment(ctx context.Context, id uuid.UUID, status PaymentStatus, eventType PaymentEventType, sessionID string, source string, providerEventID string) (*Payment, error)
}

type paymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) PaymentRepo {
	return &paymentRepo{
		db: db,
	}
}

func (r *paymentRepo) FindPaymentByID(ctx context.Context, id uuid.UUID) (*Payment, error) {
	var payment Payment

	err := r.db.WithContext(ctx).
		Preload("User").
		First(&payment, "id = ?", id).Error

	return &payment, err
}

func (r *paymentRepo) FindUserPaymentByID(ctx context.Context, userId, paymentId uuid.UUID) (*Payment, error) {
	var payment Payment

	err := r.db.WithContext(ctx).
		Preload("User").
		First(&payment, "user_id = ? AND id = ?", userId, paymentId).Error

	return &payment, err
}

func (r *paymentRepo) CreatePayment(ctx context.Context, p *Payment) (*Payment, error) {
	paymentID := uuid.New()
	eventID := uuid.New()

	newPayment := *p
	newPayment.ID = paymentID
	newPayment.CreatedAt = time.Now()
	newPayment.UpdatedAt = time.Now()

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newPayment).Error; err != nil {
			return err
		}

		newEvent := PaymentEvent{
			ID:        eventID,
			PaymentID: newPayment.ID,
			EventType: EventCreated,
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

func (r *paymentRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]Payment, error) {
	var payments []Payment
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

func (r *paymentRepo) ListPendingByUser(ctx context.Context, userID uuid.UUID) ([]Payment, error) {
	var payments []Payment
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, "pending").
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// func (r *paymentRepo) UpdatePaymentProviderObject(ctx context.Context, id uuid.UUID, providerObjectID string) (*Payment, error) {

// 	payment, err := r.updatePayment(
// 		ctx,
// 		id,
// 		"",
// 		"created",
// 		providerObjectID,
// 		"app",
// 	)

// 	return payment, err
// }

func (r *paymentRepo) UpdatePayment(
	ctx context.Context,
	id uuid.UUID,
	status PaymentStatus,
	eventType PaymentEventType,
	sessionID string,
	source string,
	providerEventID string,
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

		if sessionID != "" {
			payment.ProviderSessionID = sessionID
		}

		payment.UpdatedAt = time.Now()

		if err := tx.Save(&payment).Error; err != nil {
			return err
		}

		newEvent := PaymentEvent{
			ID:              eventID,
			PaymentID:       payment.ID,
			EventType:       eventType,
			Source:          source,
			ProviderEventID: providerEventID,
			CreatedAt:       time.Now(),
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
