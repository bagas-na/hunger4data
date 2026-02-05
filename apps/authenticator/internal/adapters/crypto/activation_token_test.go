package crypto

import (
	"testing"
)

func TestGenerateActivationToken(t *testing.T) {
	t.Run("no error returned", func(t *testing.T) {
		token, err := GenerateActivationToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token == "" {
			t.Errorf("expected non-empty token")
		}
	})

	t.Run("length is always 44 characters", func(t *testing.T) {
		token, err := GenerateActivationToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(token) != 44 {
			t.Errorf("unexpected token length: got %d, want 44", len(token))
		}
	})

	t.Run("tokens differ between calls", func(t *testing.T) {
		token1, err := GenerateActivationToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		token2, err := GenerateActivationToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token1 == token2 {
			t.Errorf("expected different tokens, got identical values")
		}
	})
}
