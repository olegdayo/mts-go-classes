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

	r.Get("/ping", connectionCheckHandler)
	r.Get("/users", s.getUsersHandler)

	r.Get("/login", s.loginHandler)
	r.Get("/verify", s.verifyHandler)

	r.Post("/registration", s.registrationHandler)

	r.Route("/user", func(r chi.Router) {
		r.Get("/{username}", s.getUserHandler)
		r.Put("/{username}", s.updateUserHandler)
		r.Delete("/{username}", s.deleteUserHandler)
	})

	return r
}

func (s *Server) Start() {
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Server running error: %s\n", err.Error())
	}
}
