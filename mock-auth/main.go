package main

import (
	"context"
	"log"
	"mock-auth/config"
	"mock-auth/server"
	"os/signal"
	"syscall"
	"time"
)

var conf *config.Config

func init() {
	conf = new(config.Config)
	err := conf.Init()
	if err != nil {
		log.Fatalf("Config error: %s", err.Error())
	}
}

func main() {
	s := server.NewServer(conf.Server.Port)
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go s.Start()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := s.Shutdown(shutdownCtx)

	if err != nil {
		log.Fatalf("Server shutdown error: %s", err.Error())
	}
}
