package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var hmacSecret []byte

// GetJWT - create and return a JWT Token
func GetJWT() (string, error) {
	jwt, err := mockJWT()

	return jwt, err
}

func mockJWT() (string, error) {
	// Create JWT - Define alg and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_id": "123456",
		"name":        "Mr. RunTask",
		"nbf":         time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign JWT
	tokenString, err := token.SignedString(hmacSecret)

	return tokenString, err

}
