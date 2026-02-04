package repo

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

//go:generate mockery --name UserRepo --inpackage
type UserRepo interface {
	CreateUser(u User) error
	GetByUsername(username string) (*User, error)
	UpdateUser(username string, user User) error
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

func (r *GORMRepository) CreateUser(u User) error {
	err := r.db.Create(&u).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr.Code == "23505" {
			return gorm.ErrDuplicatedKey
		}
		return err
	}

	return nil
}

func (r *GORMRepository) GetByUsername(username string) (*User, error) {
	var user User
	if err := r.db.Where("username = ? ", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GORMRepository) UpdateUser(username string, user User) error {

	err := r.db.Where("username = ?", username).Updates(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GORMRepository) DeleteUser(username string) error {
	var user User
	err := r.db.Delete(user, "username = ?", username).Error
	if err != nil {
		return err
	}

	return nil
}
