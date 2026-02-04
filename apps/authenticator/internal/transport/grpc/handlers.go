package handler

import (
	"authenticator/internal/service"
	"context"
	"errors"
	"fmt"
	authenticatorv1 "hunger4data/pb/authenticator"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuthService struct {
	serv service.AuthFunc
	authenticatorv1.UnimplementedAuthServiceServer
}

func NewHandService(serv service.AuthFunc) AuthService {
	return AuthService{serv: serv}
}

func (s *AuthService) Login(ctx context.Context, req *authenticatorv1.LoginRequest) (*authenticatorv1.LoginResponse, error) {

	if req.Username == "" || req.Password == "" {
		return &authenticatorv1.LoginResponse{}, status.Error(codes.InvalidArgument, "username and password are required")
	}
	username := req.Username
	password := req.Password
	token, err := s.serv.Login(username, password)
	if err != nil {
		return &authenticatorv1.LoginResponse{Token: "", Message: "Error loggin in"}, status.Error(codes.Internal, fmt.Sprintf("%s", err))
	}
	return &authenticatorv1.LoginResponse{
		Token:   token,
		Message: "Success you are logged in",
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *authenticatorv1.RegisterRequest) (*authenticatorv1.RegisterResponse, error) {
	username := req.Username
	password := req.Password
	newUser, err := s.serv.Register(username, password)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &authenticatorv1.RegisterResponse{
				User:    &authenticatorv1.User{},
				Message: "User already exists",
			}, status.Error(codes.AlreadyExists, err.Error())
		} else if errors.Is(err, service.ErrInvalidEmail) {
			return &authenticatorv1.RegisterResponse{
				User:    &authenticatorv1.User{},
				Message: "Username must be a valid email",
			}, status.Error(codes.InvalidArgument, err.Error())
		}

		return &authenticatorv1.RegisterResponse{
			User:    &authenticatorv1.User{},
			Message: "Creating user error",
		}, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authenticatorv1.RegisterResponse{
		User: &authenticatorv1.User{
			Id:       newUser.Id.String(),
			Username: newUser.Username,
			Role:     newUser.Role,
		},
		Message: "Success you are registered",
	}, nil
}
