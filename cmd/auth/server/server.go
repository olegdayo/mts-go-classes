package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

const PathToUsersConfig = "server/users.yaml"

type Server struct {
	http.Server
	Users map[string]*User
}

func NewServer(port uint16) (s *Server) {
	s = new(Server)
	s.Addr = fmt.Sprintf(":%d", port)

	s.Handler = setRouter()

	s.Users = make(map[string]*User)
	s.importUsers(PathToUsersConfig)

	return
}

func setRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/login", login)
	r.Get("/verify", verify)
	return r
}

func (s *Server) importUsers(path string) {
}

func (s *Server) Start() {
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Server running error: %s\n", err.Error())
	}
}
