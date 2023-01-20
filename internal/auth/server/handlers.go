package server

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func writeWrapper(status int, responseMessage []byte, logMessage string, w http.ResponseWriter) {
	log.Printf("%s: %d", logMessage, status)
	w.WriteHeader(status)
	_, err := w.Write(responseMessage)
	if err != nil {
		log.Printf("Failed to write: %s\n", err.Error())
	}
}

func (s *Server) registration(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		writeWrapper(http.StatusBadRequest, []byte("username and password should not be empty"), "Login and password cannot be empty, status code", w)
		return
	}

	exists, err := s.doesUserExist(username)
	if err != nil {
		writeWrapper(http.StatusInternalServerError, []byte(err.Error()), "Error during user existence check", w)
		return
	}
	if exists {
		writeWrapper(http.StatusForbidden, []byte("User already exists"), "Failed to insert because such user already exists", w)
		return
	}

	_, err = s.DBClient.Database("users").Collection("users").InsertOne(context.Background(), User{username, password})
	if err != nil {
		writeWrapper(http.StatusInternalServerError, []byte(err.Error()), "Failed to execute DB query", w)
	}

	writeWrapper(http.StatusCreated, []byte(fmt.Sprintf("Operation was successful, username: %s", username)), fmt.Sprintf("Operation was successful, username: %s", username), w)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		writeWrapper(http.StatusForbidden, []byte(fmt.Sprintf("%d, %s", http.StatusForbidden, "Forbidden")), "Cannot auth", w)
		log.Printf("Forbidden to login, status code: %d", http.StatusForbidden)
		return
	}

	user := new(User)
	err := s.DBClient.Database("users").Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			writeWrapper(http.StatusForbidden, []byte(err.Error()), "User does not exist", w)
			return
		}
		writeWrapper(http.StatusInternalServerError, []byte(err.Error()), "Error during user existence check", w)
		return
	}

	if password != user.Password {
		writeWrapper(http.StatusForbidden, []byte("Wrong password"), "Wrong password", w)
		return
	}

	writeWrapper(http.StatusOK, []byte("Logged in"), "Logged in", w)
}

func (s *Server) verify(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	users := make([]*User, 0)
	coll, err := s.DBClient.Database("users").Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		writeWrapper(http.StatusInternalServerError, []byte(err.Error()), "failed to execute users query", w)
	}

	user := new(User)
	for coll.Next(context.TODO()) {
		err = coll.Decode(user)
		if err != nil {
			writeWrapper(http.StatusInternalServerError, []byte(err.Error()), "failed to decode users", w)
		}
		users = append(users, user)
	}

	byteUsers, err := json.Marshal(users)
	if err != nil {
		writeWrapper(http.StatusInternalServerError, []byte(err.Error()), "failed to encode users", w)
	}
	writeWrapper(http.StatusOK, byteUsers, "sent users", w)
}
