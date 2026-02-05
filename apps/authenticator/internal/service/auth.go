package service

import (
	"authenticator/internal/adapters/crypto"
	"authenticator/internal/adapters/notification"
	"authenticator/internal/adapters/repo"
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

//go:generate mockery --name AuthFunc --inpackage
type AuthFunc interface {
	Login(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, username string, password string) (*repo.User, error)
}

type AuthService struct {
	repo   repo.UserRepo
	jwt    crypto.Cryptofuncs
	mailer notification.Mailer
}

func NewAuthService(repo repo.UserRepo, jwt crypto.Cryptofuncs, mailer notification.Mailer) AuthFunc {
	return &AuthService{
		repo:   repo,
		jwt:    jwt,
		mailer: mailer,
	}
}

func (s *AuthService) Login(ctx context.Context, username string, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("Needs username and password")
	}
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("Invalid username or password")
	}

	fmt.Printf("Comparing hashed-password and pasword:\nHash: %q\nPassword: %q\n", user.PasswordHash, password)

	if !s.jwt.PassCompare(password, user.PasswordHash) {
		return "", errors.New("Invalid username or password")
	}
	token, err := s.jwt.GenerateToken(user.Id)
	if err != nil {
		return "", err
	}
	return token, nil
}

var ErrInvalidEmail = errors.New("invalid email")

func (s *AuthService) Register(ctx context.Context, username string, password string) (*repo.User, error) {
	passhashed, err := s.jwt.PassHash(password)
	if err != nil {
		return nil, err
	}

	addr, err := mail.ParseAddress(username)
	if err != nil || addr.Address != username {
		return nil, ErrInvalidEmail
	}

	activationString, err := crypto.GenerateActivationToken()
	if err != nil {
		return nil, err
	}

	user := repo.User{
		Id:               uuid.New(),
		Username:         username,
		PasswordHash:     passhashed,
		Role:             "user",
		ActivationString: activationString,
		Created_At:       time.Now(),
		Updated_At:       time.Now(),
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	err = s.mailer.SendRegistrationActivationLink(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
