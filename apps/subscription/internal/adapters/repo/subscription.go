package repo

import (
	"errors"
	"fmt"
	"subscription/internal/adapters/model"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type SubscriptionRepo interface {
	CreateSubcription(u model.Subscription) error
	GetSubscriptionsByUserID(userId uuid.UUID) ([]model.Subscription, error)
	// UpdateSubscription(id uuid.UUID, subs model.Subscription) error
	DeleteSubscription(userId, subcriptionId uuid.UUID) error
}

type GORMRepository struct {
	db *gorm.DB
}

func NewSubRepo(db *gorm.DB) SubscriptionRepo {
	return &GORMRepository{
		db: db,
	}
}

func (r *GORMRepository) CreateSubcription(subs model.Subscription) error {
	newSub := model.Subscription{
		Id:          uuid.New(),
		UserID:      subs.UserID,
		CountryCode: subs.CountryCode,
		CreatedAt:   time.Now(),
	}

	err := r.db.Create(&newSub).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return gorm.ErrDuplicatedKey
		} else if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return gorm.ErrForeignKeyViolated
		}
		return err
	}

	return nil
}

func (r *GORMRepository) GetSubscriptionsByUserID(userId uuid.UUID) ([]model.Subscription, error) {
	var sub []model.Subscription
	result := r.db.Where("user_id = ? AND deleted_at IS NULL", userId).Find(&sub).Error
	if result != nil {
		return nil, result
	}
	return sub, nil
}

// func (r *GORMRepository) UpdateSubscription(id uuid.UUID, subs model.Subscription) error {
// 	err := r.db.Where("id = ?", id).Updates(subs).Error
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *GORMRepository) DeleteSubscription(userId, subcriptionId uuid.UUID) error {
	var subs model.Subscription

	if err := r.db.
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", subcriptionId, userId).
		First(&subs).
		Error; err != nil {
		return fmt.Errorf("error fetching subscription id; %w", err)
	}

	timeNow := time.Now()
	subs.DeletedAt = &timeNow

	if err := r.db.Save(&subs).Error; err != nil {
		return err
	}
	return nil
}
