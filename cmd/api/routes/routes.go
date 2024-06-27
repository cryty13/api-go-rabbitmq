package routes

import (
	"log"
	"net/http"
	"time"

	service "projeto-api/cmd/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

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

	mux.Group(func(r chi.Router) {
		r.Use(service.JWTMiddleware)
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

// Handler protegido que requer autenticação via token JWT
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Recurso protegido"))
}
