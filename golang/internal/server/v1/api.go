package v1

import (
	"net/http"

	"github.com/chipocrudos/microblog/internal/server/data"
	"github.com/go-chi/chi/v5"
)

var follow string = "follow"

func New() http.Handler {
	r := chi.NewRouter()

	ur := &UserRouter{
		Repository: &data.UserRepository{
			Data:       &data.MongoCN,
			Collection: "users",
			Follow:     follow,
		},
	}

	pr := &PostRouter{
		Repository: &data.PostRepository{
			Data:       &data.MongoCN,
			Collection: "posts",
			Follow:     follow,
		},
	}

	r.Mount("/users", ur.Routes())
	r.Mount("/posts", pr.Routes())

	return r

}
