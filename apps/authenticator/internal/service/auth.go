package service

import (
	"authenticator/internal/adapters/crypto"
	"authenticator/internal/adapters/repo"
	"errors"
	"time"
)

type AuthFunc interface {
	Login(username string, password string) (string, error)
	Register(username string, password string) error
}

type AuthService struct {
	repo repo.UserRepo
	jwt  crypto.Cryptofuncs
}

func NewAuthService(repo repo.UserRepo) AuthFunc {
	return &AuthService{repo: repo}
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
	token, err := s.jwt.GenerateToken(user.Id, user.Username, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Register(username string, password string) error {
	existingUser, _ := s.repo.GetByUsername(username)
	if existingUser != nil {
		return errors.New("Username already exists")
	}
	if username == "" || password == "" {
		return errors.New("Needs username and password")
	}
	passhashed, err := s.jwt.PassHash(password)
	if err != nil {
		return err
	}
	user := repo.Users{
		Username:   username,
		Password:   passhashed,
		Role:       "user",
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
