package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTToken(userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expira em 24 horas

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret")) // Substitua "secret" pela sua chave secreta

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
