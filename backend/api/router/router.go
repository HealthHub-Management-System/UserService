package router

import (
	"github.com/go-chi/chi/v5"

	"backend/api/resource/health"
	"backend/api/resource/users"

	_ "backend/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	// Health check
	r.Get("/livez", health.Read)

	// Swagger API documentation
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Users API
	r.Route("/api/v1", func(r chi.Router) {
		usersAPI := &users.API{}
		r.Get("/users", usersAPI.List)
		r.Post("/users", usersAPI.Create)
		r.Get("/users/{id}", usersAPI.Read)
		r.Put("/users/{id}", usersAPI.Update)
		r.Delete("/users/{id}", usersAPI.Delete)
	})

	return r
}
