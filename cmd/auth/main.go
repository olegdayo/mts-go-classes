package main

import (
	"context"
	"github.com/offluck/mts-go-classes/cmd/auth/config"
	"github.com/offluck/mts-go-classes/cmd/auth/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var conf *config.Config

func init() {
	conf = new(config.Config)
	err := conf.Init()
	if err != nil {
		log.Fatalf("Config error: %s\n", err.Error())
	}
}

func main() {
	s := server.NewServer(conf.Server.Port)
	serverClose := make(chan os.Signal)
	signal.Notify(serverClose, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go s.Start()
	log.Println("Server start")
	<-serverClose
	log.Println("Server stop")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := s.Shutdown(shutdownCtx)

	if err != nil {
		log.Fatalf("Server shutdown error: %s\n", err.Error())
	}
}
