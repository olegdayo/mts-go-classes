package main

import (
	"context"
	"github.com/offluck/mts-go-classes/internal/auth/config"
	"github.com/offluck/mts-go-classes/internal/auth/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	//"github.com/offluck/mts-go-classes/cmd/auth/config"
	//"github.com/offluck/mts-go-classes/cmd/auth/server"
)

func connectDB() {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
}

func main() {
	conf, err := config.Init()
	if err != nil {
		log.Fatalf("Cannot initialize config: %v\n", err)
	}

	s := server.NewServer(conf.Server.Port)
	serverClose := make(chan os.Signal)
	signal.Notify(serverClose, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go s.Start()
	log.Println("Server start")
	<-serverClose
	log.Println("Server stop")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = s.Shutdown(shutdownCtx)

	if err != nil {
		log.Fatalf("Server shutdown error: %s\n", err.Error())
	}
}
