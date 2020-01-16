package web

import (
	"aurum/db"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/**
@apiDefine admin Admin user
Only available to admins, the first user of the server is by default admin.
*/

/**
@api {get} /me Request user info
@apiName GetUser
@apiGroup User
@apiHeader (Authorization) {String} Authorization Users JWT Token
@apiHeaderExample {String} Authorization Example:
                Authorization: "Bearer <token>"
@apiSuccess {String} username The username of the user
@apiSuccess {String} email The E-Mail of the user
@apiSuccess {Number} role The role of the user (0 = UserDAL, 1 = Admin)
@apiSuccess {Boolean} blocked If the user is blocked
@apiSuccessExample {json} Success Response:
	{
		"username":"victor",
		"email":"victor@example.com",
		"role":0,
		"blocked": false
	}

@apiError 404 If the user does not exist (anymore).
@apiVersion 0.0.0
*/
func (e *Endpoints) getMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(*db.User)

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

/**
@api {get} /users Get all users
@apiDescription Gets all users currently registered
@apiName GetUsers
@apiGroup User
@apiPermission admin
@apiHeader (Authorization) {String} Authorization Users JWT Token
@apiHeaderExample {String} Authorization Example:
                Authorization: "Bearer <token>"
@apiVersion 0.0.0
*/
func (e *Endpoints) getUsers(w http.ResponseWriter, req *http.Request) {
	// TODO: Use query parameters

	user := req.Context().Value(contextKeyUser).(*db.User)

	if user.Role != db.AdminRoleID {
		http.Error(w, "You're not an admin!", http.StatusUnauthorized)
		return
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
