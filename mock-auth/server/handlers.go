package server

import "net/http"

func login(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("login"))
	if err != nil {
		panic(err)
	}
}

func verify(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("verify"))
	if err != nil {
		panic(err)
	}
}
