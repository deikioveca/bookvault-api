package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken			= errors.New("invalid token")
)


func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	return []byte(secret)
}


func GenerateToken(username, role string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":		role,
		"exp":		time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}


func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidSigningMethod
        }
        return getSecret(), nil
    })

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}