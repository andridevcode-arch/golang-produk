package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengubah password menjadi hash untuk disimpan di database.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// VerifyPassword memeriksa password yang dimasukkan user.
func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	return err == nil
}

// GenerateRandomPassword membuat password acak.
// Gunakan hanya jika alur aplikasi memang memerlukannya.
func GenerateRandomPassword() (string, error) {
	bytes := make([]byte, 24)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
