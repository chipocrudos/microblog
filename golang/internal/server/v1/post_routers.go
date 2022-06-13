package v1

import (
	"net/http"

	"github.com/chipocrudos/microblog/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func (pr *PostRouter) Routes() http.Handler {
	r := chi.NewRouter()

	// Public routes

	// Protected routes

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(token.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/", pr.GetAllHandler)
		r.Post("/", pr.CreateHandler)
		r.Delete("/", pr.DeleteHandler)

		r.Get("/relate", pr.GetRelatePostHandler)
	})

	return r

}
