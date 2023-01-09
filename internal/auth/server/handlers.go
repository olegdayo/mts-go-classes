package server

import (
	"fmt"
	"log"
	"net/http"
)

func checkUser(username string, password string, w http.ResponseWriter) {
	ok, err := checkAuth(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorWrapper(fmt.Sprintf("%d, %s", http.StatusForbidden, "Forbidden"), w)

		log.Printf("An error occured while checking auth: %v, status code: %d", err, http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusForbidden)
		writeErrorWrapper(fmt.Sprintf("%d, %s", http.StatusForbidden, "Forbidden"), w)

		log.Printf("Forbidden to login, you didn't pass check, status code: %d", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	writeErrorWrapper(fmt.Sprintf("Operation was successful, username: %s", login), w)

	log.Printf("Successfully logged in, status code: %s", http.StatusOK)
	return
}

func writeErrorWrapper(message string, w http.ResponseWriter) {
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Printf("Failed to write: %s\n", err.Error())
	}
}

func signIn(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeErrorWrapper("username and password should not be empty", w)

		log.Printf("Login and password cannot be empty, status code: %d", http.StatusBadRequest)
		return
	}

	checkUser(username, password, w)
}

func login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		writeErrorWrapper(fmt.Sprintf("%d, %s", http.StatusForbidden, "Forbidden"), w)

		log.Printf("Forbidden to login, status code: %d", http.StatusForbidden)
		return
	}

	checkUser(username, password, w)
}

func verify(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	writeErrorWrapper("verify", w)
}
