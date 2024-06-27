package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extrair o token JWT do cabeçalho Authorization
		tokenString := ExtractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Token não encontrado", http.StatusUnauthorized)
			return
		}

		// Parse do token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificar o método de assinatura
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Token válido, continuar com a próxima função handler
		next.ServeHTTP(w, r)
	})
}

// Função para extrair o token do cabeçalho Authorization da requisição
func ExtractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// O token está no formato "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

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
