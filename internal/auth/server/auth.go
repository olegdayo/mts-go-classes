package server

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) doesUserExist(username string) (bool, error) {
	user := new(User)
	err := s.DBClient.Database("users").Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
