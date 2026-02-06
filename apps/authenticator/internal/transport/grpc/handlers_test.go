package handler

import (
	"authenticator/internal/adapters/repo"
	mockery "authenticator/internal/mocks"
	"authenticator/internal/service"
	"context"
	"errors"
	authenticatorv1 "hunger4data/pb/authenticator"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func TestAuthHandler_Login(t *testing.T) {
	t.Run("Successful Login", func(t *testing.T) {
		mockServ := new(mockery.MockAuthFunc) //
		handler := NewHandService(mockServ)

		ctx := context.Background()
		req := &authenticatorv1.LoginRequest{
			Username: "john_doe",
			Password: "secure_password",
		}

		expectedToken := "valid-jwt-token"
		mockServ.On("Login", ctx, req.Username, req.Password).Return(expectedToken, nil)

		res, err := handler.Login(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedToken, res.Token)
		assert.Contains(t, res.Message, "Success")
		mockServ.AssertExpectations(t)
	})

	t.Run("Service Returns Error", func(t *testing.T) {
		mockServ := new(mockery.MockAuthFunc)
		handler := NewHandService(mockServ)

		req := &authenticatorv1.LoginRequest{Username: "user", Password: "pw"}

		mockServ.On("Login", context.Background(), "user", "pw").Return("", errors.New("auth failed"))

		res, err := handler.Login(context.Background(), req)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())

		assert.Equal(t, "Error logging in", res.Message)
	})
}

func TestAuthHandler_Register(t *testing.T) {
	t.Run("Successful Registration", func(t *testing.T) {
		mockServ := new(mockery.MockAuthFunc)
		handler := NewHandService(mockServ)

		req := &authenticatorv1.RegisterRequest{
			Username: "newuser@example.com",
			Password: "password123",
		}

		returnedUser := &repo.User{
			Id:       uuid.New(),
			Username: req.Username,
			Role:     "user",
		}

		mockServ.On("Register", context.Background(), req.Username, req.Password).Return(returnedUser, nil)

		res, err := handler.Register(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, "Registration complete.d An activation link has been sent to your email.", res.Message)
		assert.Equal(t, returnedUser.Id.String(), res.User.Id)
		assert.Equal(t, returnedUser.Username, res.User.Username)
	})

	t.Run("Already Exists Error Mapping", func(t *testing.T) {
		mockServ := new(mockery.MockAuthFunc)
		handler := NewHandService(mockServ)

		req := &authenticatorv1.RegisterRequest{Username: "taken@test.com", Password: "123"}

		mockServ.On("Register", context.Background(), req.Username, req.Password).
			Return(nil, gorm.ErrDuplicatedKey)

		res, err := handler.Register(context.Background(), req)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.AlreadyExists, st.Code())
		assert.Equal(t, "User already exists", res.Message)
	})

	t.Run("Invalid Email Error Mapping", func(t *testing.T) {
		mockServ := new(mockery.MockAuthFunc)
		handler := NewHandService(mockServ)

		req := &authenticatorv1.RegisterRequest{Username: "bad-email", Password: "123"}

		mockServ.On("Register", context.Background(), req.Username, req.Password).
			Return(nil, service.ErrInvalidEmail)

		res, err := handler.Register(context.Background(), req)

		st, _ := status.FromError(err)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Equal(t, "Username must be a valid email", res.Message)
	})
}
