package crypto

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Cryptofuncs interface {
	GenerateToken(user_ID uuid.UUID) (string, error)
	PassHash(pass string) (string, error)
	PassCompare(pass string, hash string) bool
}

type jwt_const struct {
	secret   []byte
	exp_time time.Duration
}

func NewJwtPass(secret string, jwtDuration time.Duration) Cryptofuncs {
	return &jwt_const{
		secret:   []byte(secret),
		exp_time: jwtDuration,
	}
}

// type authClaims struct {
// 	UserID   uuid.UUID `json:"user_id"`
// 	Username string    `json:"username"`
// 	Role     string    `json:"Role"`
// 	jwt.RegisteredClaims
// }

func (s *jwt_const) GenerateToken(user_ID uuid.UUID) (string, error) {
	// username = strings.TrimSpace(username)
	// role = strings.TrimSpace(role)
	if user_ID == uuid.Nil {
		return "", errors.New("User id must not be empty")
	}

	claims := jwt.RegisteredClaims{
		Subject:   user_ID.String(),
		Issuer:    "auth-service",
		Audience:  []string{"bookms"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.exp_time)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(s.secret)
}
