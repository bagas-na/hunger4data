package repo

import (
	"errors"

	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(u Users) error
	GetByUsername(username string) (*Users, error)
	UpdateUser(username string, user Users) error
	DeleteUser(username string) error
}

type GORMRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &GORMRepository{
		db: db,
	}
}

func (r *GORMRepository) CreateUser(u Users) error {

	err := r.db.Create(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GORMRepository) GetByUsername(username string) (*Users, error) {
	var user Users
	result := r.db.Where("username = ? ", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *GORMRepository) UpdateUser(username string, user Users) error {

	err := r.db.Where("username = ?", username).Updates(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GORMRepository) DeleteUser(username string) error {
	var user Users
	err := r.db.Delete(user, "username = ?", username).Error
	if err != nil {
		return err
	}

	return nil
}
