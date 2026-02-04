package crypto

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	secret := "change-this-super-secret-but-insecure-key"
	duration := time.Hour
	issuer := NewJwtPass(secret, duration)
	userid := uuid.New()
	t.Run("success with valid uuid", func(t *testing.T) {

		token, err := issuer.GenerateToken(userid)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
	t.Run("fail with nil uuid", func(t *testing.T) {

		token, err := issuer.GenerateToken(uuid.Nil)
		assert.Empty(t, token)
		assert.Error(t, err)
		assert.Equal(t, "User id must not be empty", err.Error())
	})
	t.Run("success with correct user id", func(t *testing.T) {

		tokenstr, err := issuer.GenerateToken(userid)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenstr)
		token, err := jwt.ParseWithClaims(tokenstr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		assert.NoError(t, err)
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		assert.True(t, ok)

		assert.Equal(t, userid.String(), claims.Subject)
		assert.Equal(t, "auth-service", claims.Issuer)
		assert.Contains(t, claims.Audience, "bookms")

		assert.WithinDuration(t, time.Now().Add(duration), claims.ExpiresAt.Time, 2*time.Second)
	})
}
