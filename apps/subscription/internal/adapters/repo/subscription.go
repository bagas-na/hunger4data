package repo

import (
	"errors"
	"subscription/internal/service"

	"gorm.io/gorm"
)

type GORMRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) GORMRepository {
	return GORMRepository{
		db: db,
	}
}

func (r *GORMRepository) CreateSubcription(u service.Subscription) error {

	err := r.db.Create(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GORMRepository) GetBySubscriptionID(id int) (*service.Subscription, error) {
	var sub service.Subscription
	result := r.db.Where("id = ? ", id).First(&sub)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &sub, nil
}

func (r *GORMRepository) UpdateSubscription(id int, subs service.Subscription) error {

	err := r.db.Where("id = ?", id).Updates(subs).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GORMRepository) DeleteUser(id int) error {
	var user service.Subscription
	err := r.db.Delete(user, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}
