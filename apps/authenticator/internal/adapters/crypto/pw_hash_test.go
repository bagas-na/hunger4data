package crypto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPwHashing(t *testing.T) {

	jwt := NewJwtPass("change-this-super-secret-but-insecure-key", time.Hour)
	assert.False(t, jwt.PassCompare("test", "not-a-hash"))
	t.Run("success with valid pass", func(t *testing.T) {
		password := "hash"
		hashed, err := jwt.PassHash(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)

		boolean := jwt.PassCompare(password, hashed)
		assert.Equal(t, boolean, true)
	})
	t.Run("fail with wrong pass", func(t *testing.T) {
		password := "hash"
		wrongpass := "pass"
		hashed, err := jwt.PassHash(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)

		boolean := jwt.PassCompare(wrongpass, hashed)
		assert.Equal(t, boolean, false)
	})

}
