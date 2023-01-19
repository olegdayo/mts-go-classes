package server

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func writeWrapper(status int, response_message string, log_message string, w http.ResponseWriter) {
	log.Printf("%s: %d", log_message, status)
	w.WriteHeader(status)
	_, err := w.Write([]byte(response_message))
	if err != nil {
		log.Printf("Failed to write: %s\n", err.Error())
	}
}

func (s *Server) registration(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		writeWrapper(http.StatusBadRequest, "username and password should not be empty", "Login and password cannot be empty, status code", w)
		return
	}

	exists, err := s.doesUserExist(username)
	if err != nil {
		writeWrapper(http.StatusInternalServerError, err.Error(), "Error during user existence check", w)
		return
	}
	if exists {
		writeWrapper(http.StatusForbidden, "User already exists", "Failed to insert because such user already exists", w)
		return
	}

	_, err = s.DBClient.Database("users").Collection("users").InsertOne(context.Background(), User{username, password})
	if err != nil {
		writeWrapper(http.StatusInternalServerError, err.Error(), "Failed to execute DB query", w)
	}

	writeWrapper(http.StatusOK, fmt.Sprintf("Operation was successful, username: %s", username), fmt.Sprintf("Operation was successful, username: %s", username), w)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		writeWrapper(http.StatusForbidden, fmt.Sprintf("%d, %s", http.StatusForbidden, "Forbidden"), "Cannot auth", w)
		log.Printf("Forbidden to login, status code: %d", http.StatusForbidden)
		return
	}

	user := new(User)
	err := s.DBClient.Database("users").Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			writeWrapper(http.StatusForbidden, err.Error(), "User does not exist", w)
			return
		}
		writeWrapper(http.StatusInternalServerError, err.Error(), "Error during user existence check", w)
		return
	}

	if password != user.Password {
		writeWrapper(http.StatusForbidden, "Wrong password", "Wrong password", w)
		return
	}

	writeWrapper(http.StatusOK, "Logged in", "Logged in", w)
}

func (s *Server) verify(w http.ResponseWriter, r *http.Request) {
	//writeWrapper(http.StatusOK, "verify", w)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	//writeWrapper(http.StatusOK, "verify", w)
}

func (s *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	//writeWrapper(http.StatusOK, "verify", w)
}
