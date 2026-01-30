package repo

import (
	"authenticator/internal/service"

	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(u service.Users) error
	GetByUsername(username string) (*service.Users, error)
	UpdateUser(username string, user service.Users) error
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

func (r *GORMRepository) CreateUser(u service.Users) error {

	err := r.db.Create(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GORMRepository) GetByUsername(username string) (*service.Users, error) {
	var user service.Users
	result := r.db.Where("username = ? ", username).First(&user)
	if result.Error != nil {
		return &service.Users{}, result.Error
	}
	return &user, nil
}

func (r *GORMRepository) UpdateUser(username string, user service.Users) error {

	err := r.db.Where("username = ?", username).Updates(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GORMRepository) DeleteUser(username string) error {
	var user service.Users
	err := r.db.Delete(user, "username = ?", username).Error
	if err != nil {
		return err
	}

	return nil
}

// func (r *MysqlRepository) GetByID(ctx context.Context, id int64) (*internal.Users, error) {
// 	query := `SELECT id, username, password FROM users WHERE id = ?`

// 	var user internal.Users
// 	err := r.db.QueryRowContext(ctx, query, id).Scan(
// 		&user.ID,
// 		&user.Username,
// 		&user.Password,
// 		&user.CreatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }
