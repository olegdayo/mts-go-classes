package server

import (
	"log"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("login"))
	if err != nil {
		log.Fatalf("Login error: %s\n", err.Error())
	}
}

func verify(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("verify"))
	if err != nil {
		log.Fatalf("Verify error: %s\n", err.Error())
	}
}
