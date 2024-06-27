package routes

import (
	"log"
	"net/http"
	"strings"
	"time"

	service "projeto-api/cmd/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var jwtKey = []byte("secret")

func SetupRoutes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(requestLogger)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// mux.Use(middleware.Heartbeat("/ping"))

	mux.Group(func(r chi.Router) {
		r.Use(JWTMiddleware)
		// Rotas protegidas
		r.Get("/protected", ProtectedHandler)
	})

	mux.Post("/", service.Teste)
	return mux
}
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v", time.Since(start))
	})
}

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

// Handler protegido que requer autenticação via token JWT
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Recurso protegido"))
}
