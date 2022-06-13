package server

import (
	"log"
	"net/http"
	"time"

	v1 "github.com/chipocrudos/microblog/internal/server/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server *http.Server
}

func New(port string) (*Server, error) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/api/v1", v1.New())

	serv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := Server{server: serv}

	return &server, nil

}

func (serv *Server) Close() error {
	return nil
}

func (serv *Server) Start() {
	log.Printf("Server running on port :%s", serv.server.Addr)
	log.Fatal(serv.server.ListenAndServe())
}
