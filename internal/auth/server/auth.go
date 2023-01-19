package server

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (s *Server) doesUserExist(username string) (int, error) {
	user := new(User)
	err := s.DBClient.Database("users").Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 200, nil
		}
		return http.StatusInternalServerError, errors.New(fmt.Sprintf("An error occured while checking auth: %v, status code: %d", err, http.StatusInternalServerError))
	}

	return http.StatusForbidden, errors.New(fmt.Sprintf("Forbidden to login, you didn't pass check, status code: %d", http.StatusForbidden))
}
