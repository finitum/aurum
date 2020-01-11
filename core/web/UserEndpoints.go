package web

import (
	"aurum/db"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Returns the user info of the currently logged in user.
func (e *Endpoints) getMe(w http.ResponseWriter, r *http.Request) {
	claims, err := e.authenticateRequest(w, r)
	if err != nil {
		return
	}

	user, err := e.conn.GetUserByName(claims.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	bytes, err := json.Marshal(db.User{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Blocked:  user.Blocked,
	})

	_, err = w.Write(bytes)
	if err != nil {
		log.Error("Couldn't write to client")
	}
}

type Range struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Returns the user info of the currently logged in user.
func (e *Endpoints) getUsers(w http.ResponseWriter, req *http.Request) {
	claims, err := e.authenticateRequest(w, req)
	if err != nil {
		return
	}

	if claims.Role != db.AdminRoleID {
		user, err := e.conn.GetUserByName(claims.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Role != db.AdminRoleID {
			http.Error(w, "You're not an admin!", http.StatusUnauthorized)
			return
		}
	}

	var r Range
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		http.Error(w, "Please yeet us a valid body", http.StatusBadRequest)
		return
	}

	// check if there's at least one entry in the range
	if r.End <= r.Start {
		http.Error(w, "Invalid Range", http.StatusBadRequest)
		return
	}

	users, err := e.conn.GetUsers(r.Start, r.End)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var strippedusers []db.User
	for _, user := range users {
		strippedusers = append(strippedusers, db.User{
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Blocked:  user.Blocked,
		})
	}

	bytes, err := json.Marshal(strippedusers)

	_, err = w.Write(bytes)
	if err != nil {
		log.Error("Couldn't write to client")
	}
}
