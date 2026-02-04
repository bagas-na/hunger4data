package repo

import (
	"subscription/internal/adapters/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepo interface {
	CreateSubcription(u model.Subscription) error
	GetBySubscriptionUserID(id uuid.UUID) ([]model.Subscription, error)
	UpdateSubscription(id uuid.UUID, subs model.Subscription) error
	DeleteSubscription(id uuid.UUID) error
}

type GORMRepository struct {
	db *gorm.DB
}

func NewSubRepo(db *gorm.DB) SubscriptionRepo {
	return &GORMRepository{
		db: db,
	}
}

func (r *GORMRepository) CreateSubcription(u model.Subscription) error {
	err := r.db.Create(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GORMRepository) GetBySubscriptionUserID(id uuid.UUID) ([]model.Subscription, error) {
	var sub []model.Subscription
	result := r.db.Where("id_user = ? ", id).Find(&sub).Error
	if result != nil {
		return nil, result
	}
	return sub, nil
}

func (r *GORMRepository) UpdateSubscription(id uuid.UUID, subs model.Subscription) error {
	err := r.db.Where("id = ?", id).Updates(subs).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GORMRepository) DeleteSubscription(id uuid.UUID) error {
	var subs model.Subscription
	err := r.db.Delete(subs, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}
