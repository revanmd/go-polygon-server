package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorSignatureInvalid)
	}

	return token, nil
}

func GetJWTSecret() []byte {
	return jwtSecret
}
