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

func (s *AuthService) Activate(ctx context.Context, req *authenticatorv1.ActivateRequest) (*authenticatorv1.ActivateResponse, error) {
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "Activation string cannot be empty")
	}

	err := s.serv.Activate(ctx, req.Key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Invalid activation string")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authenticatorv1.ActivateResponse{
		Message: "Your account has been activated successfully.",
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *authenticatorv1.LoginRequest) (*authenticatorv1.LoginResponse, error) {

	if req.Username == "" || req.Password == "" {
		return &authenticatorv1.LoginResponse{},
			status.Error(codes.InvalidArgument, "Username and Password are required")
	}

	token, err := s.serv.Login(ctx, req.Username, req.Password)
	if err != nil {
		return &authenticatorv1.LoginResponse{
				Token:   "",
				Message: "Error logging in",
			},
			status.Error(codes.Unauthenticated, fmt.Sprintf("%s", err))
	}

	return &authenticatorv1.LoginResponse{
		Token:   token,
		Message: "Success you are logged in",
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *authenticatorv1.RegisterRequest) (*authenticatorv1.RegisterResponse, error) {
	username := req.Username
	password := req.Password

	if username == "" || password == "" {
		return &authenticatorv1.RegisterResponse{
			User:    &authenticatorv1.User{},
			Message: "Username and password must not be empty",
		}, status.Error(codes.InvalidArgument, "Username and password must not be empty")
	}

	newUser, err := s.serv.Register(ctx, username, password)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &authenticatorv1.RegisterResponse{
				User:    &authenticatorv1.User{},
				Message: "User already exists",
			}, status.Error(codes.AlreadyExists, "User already exists")
		} else if errors.Is(err, service.ErrInvalidEmail) {
			return &authenticatorv1.RegisterResponse{
				User:    &authenticatorv1.User{},
				Message: "Username must be a valid email",
			}, status.Error(codes.InvalidArgument, "Username must be a valid email")
		}

		return &authenticatorv1.RegisterResponse{
			User:    &authenticatorv1.User{},
			Message: "Creating user error",
		}, status.Error(codes.Internal, err.Error())
	}

	return &authenticatorv1.RegisterResponse{
		User: &authenticatorv1.User{
			Id:       newUser.Id.String(),
			Username: newUser.Username,
			Role:     newUser.Role,
		},
		Message: "Registration complete.d An activation link has been sent to your email.",
	}, nil
}
