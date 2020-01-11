package web

import (
	"aurum/db"
	"aurum/hash"
	"aurum/jwt"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// Creates a user based on the user in the json body
func (e *Endpoints) signup(w http.ResponseWriter, r *http.Request) {

	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Please yeet us a valid json body", http.StatusBadRequest)
		return
	}

	if u.Username == "" || u.Password == "" || u.Email == "" || u.Role != db.UserRoleID {
		http.Error(w, "Please yeet us a valid json body", http.StatusBadRequest)
		return
	}

	// If the signed up user is the only user, make this user admin
	number, err := e.conn.CountUsers()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if number == 0 {
		log.Infof("A user with username \"%s\" has signed up. This is the first user and will get admin privileges.", u.Username)
		u.Role = db.AdminRoleID
	}

	// Actually add the user to the db
	err = e.conn.CreateUser(u)
	if err != nil {
		// look if the error was caused by the username already existing
		// TODO: This error is SQL specific so should not be handled here probably
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed:") {
			http.Error(w, "Username already chosen", http.StatusConflict)
		} else {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

// uses the refresh token to return a new login token
func (e *Endpoints) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please specify a body", http.StatusBadRequest)
		return
	}

	var t jwt.TokenPair
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Please yeet us a valid json body", http.StatusBadRequest)
		return
	}

	c, err := jwt.VerifyJWT(t.RefreshToken, e.config)
	if err != nil || !c.Refresh {
		http.Error(w, "Please specify a valid refresh token", http.StatusBadRequest)
		return
	}

	user, err := e.conn.GetUserByName(c.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	token, err := jwt.GenerateJWT(&user, false, e.config)
	if err != nil {
		http.Error(w, "Couldn't generate login token", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(jwt.TokenPair{
		LoginToken: token,
	})

	_, err = w.Write(bytes)

	if err != nil {
		log.Error("Couldn't write to client")
	}

	return
}

// Expects a JSON body with the user object
// returns the refresh and login token
func (e *Endpoints) login(w http.ResponseWriter, r *http.Request) {

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
	if !hash.CheckPasswordHash(u.Password, user.Password) {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}

	// Generate a JWT pair (login + refresh)
	token, err := jwt.GenerateJWTPair(&user, e.config)
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
		log.Error("Error writing response to client")
	}

	return
}

// Expects a JSON body with the user object with a new password
// returns a 200 status on success
func (e *Endpoints) changePassword(w http.ResponseWriter, r *http.Request) {

	claims, err := e.authenticateRequest(w, r)
	if err != nil {
		return
	}

	// Decode user struct and check if anything is invalid.
	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil || len(u.Password) == 0 {
		http.Error(w, "Please yeet us a valid body", http.StatusBadRequest)
		return
	}

	password := u.Password

	user, err := e.conn.GetUserByName(claims.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	passwordhash, err := hash.HashPassword(password)

	user.Password = passwordhash
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = e.conn.UpdateUser(user)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	return
}
