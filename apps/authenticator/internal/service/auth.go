package service

import (
	"context"
	authenticatorv1 "hunger4data/pb/authenticator"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepo interface {
	CreateUser(u Users) error
	GetByUsername(username string) (*Users, error)
	UpdateUser(username string, user Users) error
	DeleteUser(username string) error
}

type crypto interface {
	GenerateToken(user_ID uuid.UUID, username string, role string) (string, error)
	PassHash(pass string) (string, error)
	PassCompare(pass string, hash string) bool
}

type AuthService struct {
	repo UserRepo
	authenticatorv1.UnimplementedAuthServiceServer
	jwt crypto
}

func NewAuthService(repo UserRepo, jwt crypto) AuthService {
	return AuthService{repo: repo, jwt: jwt}
}

func (s *AuthService) Login(ctx context.Context, req *authenticatorv1.LoginRequest) (*authenticatorv1.LoginResponse, error) {

	if req.Username == "" || req.Password == "" {
		return &authenticatorv1.LoginResponse{}, status.Error(codes.InvalidArgument, "username and password are required")
	}
	username := req.Username
	password := req.Password

	User, err := s.repo.GetByUsername(username)
	if err != nil {
		return &authenticatorv1.LoginResponse{Token: "", Message: "Error getting user is the user registered?"}, status.Error(codes.Internal, "invalid username or password")
	}
	if s.jwt.PassCompare(User.Password, password) {
		return &authenticatorv1.LoginResponse{Token: "", Message: "Passwor error please re enter your password"}, status.Error(codes.NotFound, "invalid username or password")
	}
	token, err := s.jwt.GenerateToken(User.Id, User.Username, User.Role)
	if err != nil {
		return &authenticatorv1.LoginResponse{Token: "", Message: "Error generating token"}, status.Error(codes.Internal, "Error generating token")
	}
	return &authenticatorv1.LoginResponse{
		Token:   token,
		Message: "Success you are logged in",
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *authenticatorv1.RegisterRequest) (*authenticatorv1.RegisterResponse, error) {

	if req.Username == "" || req.Password == "" {
		return &authenticatorv1.RegisterResponse{User: &authenticatorv1.User{}, Message: "username and password are required"}, status.Error(codes.InvalidArgument, "username and password are required")
	}
	username := req.Username
	password := req.Password
	existingUser, _ := s.repo.GetByUsername(username)
	if existingUser != nil {
		return &authenticatorv1.RegisterResponse{User: &authenticatorv1.User{}, Message: "username already exists"}, status.Error(codes.AlreadyExists, "username already exists")
	}

	PasswordHashed, err := s.jwt.PassHash(password)
	if err != nil {
		return &authenticatorv1.RegisterResponse{User: &authenticatorv1.User{}, Message: "Hashing error"}, status.Error(codes.Internal, "Creating user error")
	}

	user := Users{
		Username: username,
		Password: PasswordHashed,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return &authenticatorv1.RegisterResponse{User: &authenticatorv1.User{}, Message: "Creating user error"}, status.Error(codes.AlreadyExists, "failed at creating user")
	}

	return &authenticatorv1.RegisterResponse{
		User:    &authenticatorv1.User{},
		Message: "Success you are registered",
	}, nil
}
