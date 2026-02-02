package service

import (
	proto "authenticator/proto/authenticator"
	"context"

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
	GenerateToken(user_ID int, username string, role string) (string, error)
	PassHash(pass string) (string, error)
	PassCompare(pass string, hash string) bool
}

type AuthService struct {
	repo UserRepo
	proto.UnimplementedAuthServiceServer
	jwt crypto
}

func NewAuthService(repo UserRepo, jwt crypto) AuthService {
	return AuthService{repo: repo, jwt: jwt}
}

func (s *AuthService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {

	if req.Username == "" || req.Password == "" {
		return &proto.LoginResponse{}, status.Error(codes.InvalidArgument, "username and password are required")
	}
	username := req.Username
	password := req.Password

	User, err := s.repo.GetByUsername(username)
	if err != nil {
		return &proto.LoginResponse{Token: "", Message: "Error getting user is the user registered?"}, status.Error(codes.Internal, "invalid username or password")
	}
	if s.jwt.PassCompare(User.Password, password) {
		return &proto.LoginResponse{Token: "", Message: "Passwor error please re enter your password"}, status.Error(codes.NotFound, "invalid username or password")
	}
	token, err := s.jwt.GenerateToken(User.Id, User.Username, User.Role)
	if err != nil {
		return &proto.LoginResponse{Token: "", Message: "Error generating token"}, status.Error(codes.Internal, "Error generating token")
	}
	return &proto.LoginResponse{
		Token:   token,
		Message: "Success you are logged in",
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	if req.Username == "" || req.Password == "" {
		return &proto.RegisterResponse{User: &proto.User{}, Message: "username and password are required"}, status.Error(codes.InvalidArgument, "username and password are required")
	}
	username := req.Username
	password := req.Password
	existingUser, _ := s.repo.GetByUsername(username)
	if existingUser != nil {
		return &proto.RegisterResponse{User: &proto.User{}, Message: "username already exists"}, status.Error(codes.AlreadyExists, "username already exists")
	}

	PasswordHashed, err := s.jwt.PassHash(password)
	if err != nil {
		return &proto.RegisterResponse{User: &proto.User{}, Message: "Hashing error"}, status.Error(codes.Internal, "Creating user error")
	}

	user := Users{
		Username: username,
		Password: PasswordHashed,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return &proto.RegisterResponse{User: &proto.User{}, Message: "Creating user error"}, status.Error(codes.AlreadyExists, "failed at creating user")
	}

	return &proto.RegisterResponse{
		User:    &proto.User{},
		Message: "Success you are registered",
	}, nil
}
