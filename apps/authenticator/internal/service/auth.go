package service

import (
	"authenticator/internal/adapters/crypto"
	"authenticator/internal/adapters/repo"
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type AuthFunc interface {
	Login(username string, password string) (string, error)
	Register(username string, password string) (*repo.User, error)
}

type AuthService struct {
	repo repo.UserRepo
	jwt  crypto.Cryptofuncs
}

func NewAuthService(repo repo.UserRepo, jwt crypto.Cryptofuncs) AuthFunc {
	return &AuthService{
		repo: repo,
		jwt:  jwt,
	}
}

func (s *AuthService) Login(username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("Needs username and password")
	}
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("Cannot find username")
	}
	if s.jwt.PassCompare(user.Password, password) {
		return "", errors.New("Wrong password")
	}
	token, err := s.jwt.GenerateToken(user.Id)
	if err != nil {
		return "", err
	}
	return token, nil
}

var ErrInvalidEmail = errors.New("invalid email")

func (s *AuthService) Register(username string, password string) (*repo.User, error) {
	passhashed, err := s.jwt.PassHash(password)
	if err != nil {
		return nil, err
	}

	addr, err := mail.ParseAddress(username)
	if err != nil || addr.Address != username {
		return nil, ErrInvalidEmail
	}

	user := repo.User{
		Id:         uuid.New(),
		Username:   username,
		Password:   passhashed,
		Role:       "user",
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
