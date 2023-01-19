package main

import (
	"context"
	"github.com/offluck/mts-go-classes/internal/auth/config"
	"github.com/offluck/mts-go-classes/internal/auth/server"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	conf, err := config.Init()
	if err != nil {
		log.Fatalf("Cannot initialize config: %v\n", err)
	}

	client, err := connectDB(conf.DBURI)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	log.Println(client)

	dbNames, err := client.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Failed to list database names: %v", err)
	}
	log.Println(dbNames)
	
	s := server.NewServer(conf.Server.Port, client)
	serverClose := make(chan os.Signal)
	signal.Notify(serverClose, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go s.Start()
	log.Println("Server start")
	<-serverClose
	log.Println("Server stop")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.Shutdown(shutdownCtx)

	if err != nil {
		log.Fatalf("Server shutdown error: %s\n", err.Error())
	}
}
