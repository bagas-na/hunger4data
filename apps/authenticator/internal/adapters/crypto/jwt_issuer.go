package crypto

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwt_const struct {
	secret   []byte
	exp_time time.Duration
}

func NewJwt(secret string) jwt_const {
	return jwt_const{secret: []byte(secret)}
}

type authClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"Role"`
	jwt.RegisteredClaims
}

func (s *jwt_const) GenerateToken(user_ID int, username string, role string) (string, error) {
	user_ID_stringfy := strconv.Itoa(user_ID)
	username = strings.TrimSpace(username)
	role = strings.TrimSpace(role)
	if user_ID_stringfy == "" || username == "" || role == "" {
		return "", errors.New("Error one of the fields is empty")
	}

	claim := authClaims{
		UserID:   user_ID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.exp_time)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user_ID_stringfy,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(s.secret)
}
