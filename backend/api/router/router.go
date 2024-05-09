package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/wader/gormstore/v2"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"backend/api/resource/health"
	"backend/api/resource/users"
	"backend/api/router/middleware"

	_ "backend/docs" // Swagger API documentation

	httpSwagger "github.com/swaggo/http-swagger"
)

func New(l *zerolog.Logger, db *gorm.DB, v *validator.Validate, s *gormstore.Store) *chi.Mux {
	r := chi.NewRouter()

	loggerMiddleware := middleware.NewLogger(l)

	// Health check
	r.Get("/livez", health.Read)

	// Swagger API documentation
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Users API
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
			Debug:            true,
		}))
		r.Use(middleware.ContentTypeJSON)
		r.Use(loggerMiddleware)

		usersAPI := users.New(l, db, v, s)
		r.Get("/users", usersAPI.List)
		r.Post("/users", usersAPI.Create)
		r.Get("/users/{id}", usersAPI.Read)
		r.Put("/users/{id}", usersAPI.Update)
		r.Delete("/users/{id}", usersAPI.Delete)
		r.Post("/users/login", usersAPI.Login)
		r.Post("/users/logout", usersAPI.Logout)
	})

	return r
}
