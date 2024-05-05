package router

import (
	"github.com/go-chi/chi"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"backend/api/resource/health"
	"backend/api/resource/users"
	"backend/api/router/middleware"

	_ "backend/docs" // Swagger API documentation

	httpSwagger "github.com/swaggo/http-swagger"
)

func New(l *zerolog.Logger, db *gorm.DB, v *validator.Validate) *chi.Mux {
	r := chi.NewRouter()
	loggerMiddleware := middleware.NewLogger(l)

	// Health check
	r.Get("/livez", health.Read)

	// Swagger API documentation
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Users API
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJSON)
		r.Use(loggerMiddleware)

		usersAPI := users.New(l, db, v)
		r.Get("/users", usersAPI.List)
		r.Post("/users", usersAPI.Create)
		r.Get("/users/{id}", usersAPI.Read)
		r.Put("/users/{id}", usersAPI.Update)
		r.Delete("/users/{id}", usersAPI.Delete)
	})

	return r
}
