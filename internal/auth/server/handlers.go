package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func writeWrapper(status int, message string, w http.ResponseWriter) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Printf("Failed to write: %s\n", err.Error())
	}
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		writeWrapper(http.StatusBadRequest, "username and password should not be empty", w)

		log.Printf("Login and password cannot be empty, status code: %d", http.StatusBadRequest)
		return
	}

	status, err := s.doesUserExist(username)
	if err != nil {
		writeWrapper(status, err.Error(), w)
		return
	}

	_, err = s.DBClient.Database("users").Collection("users").InsertOne(context.Background(), User{username, password})

	if err != nil {
		writeWrapper(http.StatusInternalServerError, err.Error(), w)
	}

	writeWrapper(http.StatusOK, fmt.Sprintf("Operation was successful, username: %s", username), w)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	//username, password, ok := r.BasicAuth()
	//if !ok {
	//	writeWrapper(http.StatusForbidden, fmt.Sprintf("%d, %s", http.StatusForbidden, "Forbidden"), w)
	//
	//	log.Printf("Forbidden to login, status code: %d", http.StatusForbidden)
	//	return
	//}
	//
	//s.checkUser(username, password, w)
}

func (s *Server) verify(w http.ResponseWriter, r *http.Request) {
	writeWrapper(http.StatusOK, "verify", w)
}
