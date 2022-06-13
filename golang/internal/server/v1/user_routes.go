package v1

import (
	"net/http"

	"github.com/chipocrudos/microblog/pkg/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Post("/", ur.CreateHandler)
		r.Post("/auth", ur.LoginHandler)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(token.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/", ur.GetAllHandler)

		r.Put("/", ur.UpdateHandler)
		r.Get("/me", ur.GetMeHandler)

		r.Get("/follower", ur.CreateRelationHandler)
		r.Delete("/follower", ur.DeleteRelationHandler)

		r.Get("/{id}", ur.GetUserByIdHandler)

	})

	return r
}
