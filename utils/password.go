package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword mengubah password asli menjadi hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword membandingkan password inputan user dengan hash di database
// Mengembalikan true jika cocok, false jika salah
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}