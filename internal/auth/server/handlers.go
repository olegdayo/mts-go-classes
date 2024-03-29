package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/offluck/mts-go-classes/internal/auth/encryption"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func writeWrapper(status int, message string, w http.ResponseWriter) {
	log.Printf("%s: %d", message, status)
	w.WriteHeader(status)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Printf("Failed to write: %s\n", err.Error())
	}
}

func makeMessageResponse(message string, err error) string {
	if err != nil {
		return fmt.Sprintf(`{"error": "%s: %v"}`, message, err)
	}
	return fmt.Sprintf(`{"success": "%s"}`, message)
}

func connectionCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Printf("Failed to write test response")
	}
}

func (s *Server) registrationHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := encryption.MakeHash(r.FormValue("password"))

	if username == "" || password == "" {
		writeWrapper(
			http.StatusBadRequest,
			makeMessageResponse("login and password cannot be empty", errors.New("empty data")),
			w,
		)
		return
	}

	exists, err := s.doesUserExist(username)
	if err != nil {
		writeWrapper(
			http.StatusInternalServerError,
			makeMessageResponse("error during user existence check", err),
			w,
		)
		return
	}
	if exists {
		writeWrapper(
			http.StatusForbidden,
			makeMessageResponse("user already exists", errors.New("user duplicate")),
			w,
		)
		return
	}

	_, err = s.DBClient.Database("users").Collection("users").InsertOne(
		context.Background(),
		User{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		writeWrapper(
			http.StatusInternalServerError,
			makeMessageResponse("failed to execute DB query", err),
			w,
		)
	}

	writeWrapper(
		http.StatusCreated,
		makeMessageResponse(fmt.Sprintf("operation was successful, username: %s", username), nil),
		w,
	)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		writeWrapper(
			http.StatusForbidden,
			makeMessageResponse("cannot use basic auth", errors.New("basic auth error")),
			w,
		)
		return
	}

	user := new(User)
	err := s.DBClient.Database("users").Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			writeWrapper(
				http.StatusForbidden,
				makeMessageResponse("user does not exist", err),
				w,
			)
			return
		}

		writeWrapper(
			http.StatusInternalServerError,
			makeMessageResponse("error during user existence check", err),
			w,
		)
		return
	}

	if encryption.MakeHash(password) != user.Password {
		writeWrapper(
			http.StatusForbidden,
			makeMessageResponse("wrong password", errors.New("wrong password")),
			w,
		)
		return
	}

	writeWrapper(
		http.StatusOK,
		makeMessageResponse("successfully logged in", nil),
		w,
	)
}

func (s *Server) verifyHandler(w http.ResponseWriter, r *http.Request) {
	s.loginHandler(w, r)
}

func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	log.Printf("Username: %s", username)

	user := new(User)
	err := s.DBClient.Database("users").Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeWrapper(
				http.StatusForbidden,
				makeMessageResponse("user not found", errors.New("no such user")),
				w,
			)
			return
		}

		writeWrapper(
			http.StatusInternalServerError,
			makeMessageResponse("error during decoding", err),
			w,
		)
		return
	}

	byteUser, err := json.Marshal(user)
	if err != nil {
		writeWrapper(
			http.StatusInternalServerError,
			makeMessageResponse("error during decoding", err),
			w,
		)
		return
	}

	writeWrapper(
		http.StatusOK,
		string(byteUser),
		w,
	)
}

func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	oldUsername := chi.URLParam(r, "username")
	newUsername := r.FormValue("username")
	newPassword := encryption.MakeHash(r.FormValue("password"))

	s.DBClient.Database("users").Collection("users").FindOneAndReplace(
		context.Background(),
		bson.M{"username": oldUsername},
		bson.M{
			"username": newUsername,
			"password": newPassword,
		},
	)

	writeWrapper(
		http.StatusOK,
		makeMessageResponse(fmt.Sprintf("successfully updated user: %s", oldUsername), nil),
		w,
	)
}

func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	log.Printf("Username: %s", username)
	s.DBClient.Database("users").Collection("users").FindOneAndDelete(context.Background(), bson.M{"username": username})
	writeWrapper(
		http.StatusNoContent,
		makeMessageResponse(fmt.Sprintf("successfully deleted the user: %s", username), nil),
		w,
	)
}

func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := make([]*User, 0)
	coll, err := s.DBClient.Database("users").Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		writeWrapper(
			http.StatusInternalServerError,
			makeMessageResponse("failed to execute users query", err),
			w,
		)
	}

	for coll.Next(context.TODO()) {
		user := new(User)
		err = coll.Decode(user)
		if err != nil {
			writeWrapper(
				http.StatusInternalServerError,
				makeMessageResponse("failed to decode users", err),
				w,
			)
		}
		users = append(users, user)
	}

	byteUsers, err := json.Marshal(users)
	if err != nil {
		writeWrapper(
			http.StatusInternalServerError,
			makeMessageResponse("failed to encode users", err),
			w,
		)
	}

	writeWrapper(
		http.StatusOK,
		string(byteUsers),
		w,
	)
}
