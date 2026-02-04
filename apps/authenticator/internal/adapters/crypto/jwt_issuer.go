package crypto

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Cryptofuncs interface {
	GenerateToken(user_ID uuid.UUID, username string, role string) (string, error)
	PassHash(pass string) (string, error)
	PassCompare(pass string, hash string) bool
}

type jwt_const struct {
	secret   []byte
	exp_time time.Duration
}

func NewJwtPass(secret string) Cryptofuncs {
	return &jwt_const{secret: []byte(secret)}
}

type authClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"Role"`
	jwt.RegisteredClaims
}

func (s *jwt_const) GenerateToken(user_ID uuid.UUID, username string, role string) (string, error) {
	username = strings.TrimSpace(username)
	role = strings.TrimSpace(role)
	if user_ID == uuid.Nil || username == "" || role == "" {
		return "", errors.New("Error one of the fields is empty")
	}

	claim := authClaims{
		UserID:   user_ID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.exp_time)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(s.secret)
}
