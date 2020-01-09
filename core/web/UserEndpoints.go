package web

import (
	"aurum/config"
	"aurum/db"
	"aurum/util"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Endpoints struct {
	conn   db.Connection
	config *config.Config
}

// Creates a user based on the user in the json body
func (e Endpoints) signup(w http.ResponseWriter, r *http.Request) {

	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Please yeet us a valid json body", http.StatusBadRequest)
		return
	}

	if u.Username == "" || u.Password == "" || u.Email == "" || u.Role != 0 {
		http.Error(w, "Please yeet us a valid json body", http.StatusBadRequest)
		return
	}

	// Actually add the user to the db
	err := e.conn.CreateUser(u)
	if err != nil {
		// look if the error was caused by the username already existing
		// TODO: This error is SQL specific so should not be handled here probably
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed:") {
			http.Error(w, "Server error", http.StatusConflict)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

// uses the refresh token to return a new login token
func (e Endpoints) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please specify a body", http.StatusBadRequest)
		return
	}

	return
}

// Expects a JSON body with the user object
// returns the refresh and login token
func (e Endpoints) login(w http.ResponseWriter, r *http.Request) {

	// Decode user struct and check if anything is invalid.
	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil || len(u.Username) == 0 || len(u.Password) == 0 {
		http.Error(w, "Please yeet us a valid body", http.StatusBadRequest)
		return
	}

	user, err := e.conn.GetUserByName(u.Username)
	if err != nil {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}

	// Invalid password
	if !util.CheckPasswordHash(u.Password, user.Password) {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}

	// Generate a JWT pair (login + refresh)
	token, err := GenerateJWTPair(&user, e.config)
	if err != nil {
		http.Error(w, "Error in JWT generation", http.StatusInternalServerError)
		return
	}

	// Convert the tokens into json bytes
	bytes, err := json.Marshal(token)
	if err != nil {
		http.Error(w, "Error in serializing JWT", http.StatusInternalServerError)
		return
	}

	// push them to the client
	if _, err = w.Write(bytes); err != nil {
		log.Fatalf("Error writing response to client")
	}

	return
}
