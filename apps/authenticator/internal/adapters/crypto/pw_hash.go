package crypto

import "golang.org/x/crypto/bcrypt"

func (s *jwt_const) PassHash(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	return string(hashed), err
}

func (s *jwt_const) PassCompare(pass string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return false
	}
	return true
}
