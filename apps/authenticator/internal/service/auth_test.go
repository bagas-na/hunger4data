package service

import (
	"authenticator/internal/adapters/repo"
	mockery "authenticator/internal/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Login(t *testing.T) {

	t.Run("Successful Login", func(t *testing.T) {
		mockRepo := new(mockery.MockUserRepo)
		mockJwt := new(mockery.MockCryptofuncs)
		mockMailer := new(mockery.MockMailer)
		service := NewAuthService(mockRepo, mockJwt, mockMailer)
		ctx := context.Background()
		username := "johndoe"
		password := "secret123"
		expectedToken := "fake-jwt-token"
		user := &repo.User{Id: uuid.New(), Username: "johndoe", PasswordHash: "hashed_password"}
		mockRepo.On("GetByUsername", username).Return(user, nil)
		mockJwt.On("PassCompare", password, user.PasswordHash).Return(true)
		mockJwt.On("GenerateToken", user.Id).Return(expectedToken, nil)

		token, err := service.Login(ctx, username, password)

		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
		mockRepo.AssertExpectations(t)
		mockJwt.AssertExpectations(t)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		mockRepo := new(mockery.MockUserRepo)
		mockJwt := new(mockery.MockCryptofuncs)
		mockMailer := new(mockery.MockMailer)
		service := NewAuthService(mockRepo, mockJwt, mockMailer)
		ctx := context.Background()
		username := "johndoe"
		password := "wrong_pass"
		user := &repo.User{Username: username, PasswordHash: "hashed_password"}

		mockRepo.On("GetByUsername", username).Return(user, nil)
		mockJwt.On("PassCompare", password, "hashed_password").Return(false)

		token, err := service.Login(ctx, username, password)

		assert.Error(t, err)
		assert.Equal(t, "Invalid username or password", err.Error())
		assert.Empty(t, token)
	})

	t.Run("Invalid username or password", func(t *testing.T) {
		mockRepo := new(mockery.MockUserRepo)
		mockJwt := new(mockery.MockCryptofuncs)
		mockMailer := new(mockery.MockMailer)
		service := NewAuthService(mockRepo, mockJwt, mockMailer)
		ctx := context.Background()
		username := "michael"
		password := "pass"
		mockRepo.On("GetByUsername", username).Return(nil, errors.New("db error"))

		token, err := service.Login(ctx, username, password)

		assert.Error(t, err)
		assert.Equal(t, "Invalid username or password", err.Error())
		assert.Empty(t, token)

		mockJwt.AssertNotCalled(t, "PassCompare", password, password)
		mockJwt.AssertNotCalled(t, "GenerateToken", "test")
	})

	t.Run("Needs username and password", func(t *testing.T) {
		mockRepo := new(mockery.MockUserRepo)
		mockJwt := new(mockery.MockCryptofuncs)
		mockMailer := new(mockery.MockMailer)
		service := NewAuthService(mockRepo, mockJwt, mockMailer)
		ctx := context.Background()
		username := ""
		password := ""

		token, err := service.Login(ctx, username, password)

		assert.Error(t, err)
		assert.Equal(t, "Needs username and password", err.Error())
		assert.Empty(t, token)

		mockRepo.AssertNotCalled(t, "GetByUsername", username)
		mockJwt.AssertNotCalled(t, "PassCompare", password, password)
		mockJwt.AssertNotCalled(t, "GenerateToken", "test")
	})
}

func TestAuthService_Register(t *testing.T) {

	t.Run("Successful Registration", func(t *testing.T) {
		mockRepo := new(mockery.MockUserRepo)
		mockJwt := new(mockery.MockCryptofuncs)
		mockMailer := new(mockery.MockMailer)
		service := NewAuthService(mockRepo, mockJwt, mockMailer)
		ctx := context.Background()
		email := "test@example.com"
		rawPass := "password123"
		hashedPass := "hashed_result"

		mockJwt.On("PassHash", rawPass).Return(hashedPass, nil)

		mockRepo.On("CreateUser", mock.MatchedBy(func(u repo.User) bool {
			return u.Username == email &&
				u.PasswordHash == hashedPass &&
				u.Role == "user"
		})).Return(nil)

		mockMailer.On("SendRegistrationActivationLink", mock.Anything, mock.AnythingOfType("*repo.User")).Return(nil)

		user, err := service.Register(ctx, email, rawPass)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Username)
		assert.Equal(t, "user", user.Role)
		mockRepo.AssertExpectations(t)
		mockMailer.AssertExpectations(t)
	})
	t.Run("Invalid Email Format", func(t *testing.T) {
		mockRepo := new(mockery.MockUserRepo)
		mockJwt := new(mockery.MockCryptofuncs)
		mockMailer := new(mockery.MockMailer)
		service := NewAuthService(mockRepo, mockJwt, mockMailer)
		ctx := context.Background()
		invalidEmail := "!!!"

		mockRepo.ExpectedCalls = nil
		mockJwt.ExpectedCalls = nil

		mockJwt.On("PassHash", "pass").Return("hashed_pass", nil)

		user, err := service.Register(ctx, invalidEmail, "pass")

		assert.ErrorIs(t, err, ErrInvalidEmail)
		assert.Nil(t, user)
		mockRepo.AssertNotCalled(t, "CreateUser", mock.Anything)
	})
}
