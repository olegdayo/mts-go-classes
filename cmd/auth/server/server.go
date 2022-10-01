package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

const PathToUsersConfig = "server/users.yaml"

type Server struct {
	http.Server
	Users map[string]*User `yaml:"users"`
}

func NewServer(port uint16) (s *Server) {
	s = new(Server)
	s.Addr = fmt.Sprintf(":%d", port)

	s.Handler = setRouter()

	s.Users = make(map[string]*User)
	err := s.importUsers(PathToUsersConfig)
	if err != nil {
		s.Users = make(map[string]*User)
	}
	fmt.Println(s.Users)

	return
}

func setRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/login", login)
	r.Get("/verify", verify)
	return r
}

func (s *Server) importUsers(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("cannot open users file: %s\n", err.Error())
		return err
	}

	err = yaml.Unmarshal(bytes, s)
	if err != nil {
		log.Printf("cannot import users: %s\n", err.Error())
	}

	return nil
}

func (s *Server) Start() {
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Server running error: %s\n", err.Error())
	}
}
