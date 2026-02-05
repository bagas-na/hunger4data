package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateActivationToken generates a cryptographically secure
// Base64URL-encoded token of 32 random bytes (256 bits) with padding
// for consistent length (44 chars).
func GenerateActivationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
