package server

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	http.Server
	DBClient *mongo.Client
}

func NewServer(port uint16, dbClient *mongo.Client) (s *Server) {
	s = new(Server)
	s.Addr = fmt.Sprintf(":%d", port)
	s.Handler = s.setRouter()
	s.DBClient = dbClient
	return
}

func (s *Server) setRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/users", s.getUsers)
	r.Route("/user", func(r chi.Router) {
		r.Get("/{username}", s.getUser)
		r.Put("/{username}", s.updateUser)
		r.Delete("/{username}", s.deleteUser)
	})

	r.Get("/login", s.login)
	r.Get("/verify", s.verify)

	r.Post("/registration", s.registration)

	return r
}

func (s *Server) Start() {
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Server running error: %s\n", err.Error())
	}
}
