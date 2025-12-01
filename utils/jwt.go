// package utils

// import (
// 	"os"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/google/uuid"
// )

// // JWTClaims adalah isi data yang disisipkan di dalam token
// type JWTClaims struct {
// 	UserID   uuid.UUID `json:"user_id"`
// 	Role     string    `json:"role"`
// 	jwt.RegisteredClaims
// }

// // GenerateToken membuat token baru setelah login sukses
// func GenerateToken(userID uuid.UUID, roleName string) (string, error) {
// 	// Ambil Secret Key dari .env
// 	secret := os.Getenv("JWT_SECRET")
	
// 	// Set waktu kadaluarsa token (misal: 24 jam)
// 	expirationTime := time.Now().Add(24 * time.Hour)

// 	claims := &JWTClaims{
// 		UserID: userID,
// 		Role:   roleName,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		},
// 	}

// 	// Buat token dengan algoritma HS256
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
// 	// Tandatangani token dengan secret key
// 	tokenString, err := token.SignedString([]byte(secret))
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// // ValidateToken memvalidasi token string dan mengembalikan claims jika valid
// func ValidateToken(tokenString string) (*JWTClaims, error) {
// 	secret := os.Getenv("JWT_SECRET")

// 	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(secret), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
// 		return claims, nil
// 	}

// 	return nil, err
// }

package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims adalah isi data yang disisipkan di dalam token
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken membuat token baru setelah login sukses
func GenerateToken(userID uuid.UUID, roleName string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		UserID: userID,
		Role:   roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken memvalidasi token string dan mengembalikan claims jika valid
func ValidateToken(tokenString string) (*JWTClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// SECURITY CHECK: Pastikan metode signing adalah HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}