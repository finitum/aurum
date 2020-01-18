package web

import (
	"aurum/db"
	"aurum/hash"
	"aurum/passwords"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/**
@apiDefine admin Admin user
Only available to admins, the first user of the server is by default admin.
*/

/**
@apiDefine UserObjectSuccess
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
 */

/**
@apiDefine UserObjectParam
@apiParam {String} username The username of the user
@apiParam {String} email The E-Mail of the user
@apiParam {String} password The password of the user
@apiParam {Number} role The role of the user (0 = UserDAL, 1 = Admin)
@apiParam {Boolean} blocked If the user is blocked
@apiParamExample {json} Success Response:
	{
		"username":"victor",
		"email":"victor@example.com",
  		"password": "hunter2",
		"role":0,
		"blocked": false
	}
*/

/**
@apiDefine AuthHeader
@apiHeader (Authorization) {String} Authorization User's JWT Token
@apiHeaderExample {String} Authorization Example:
                Authorization: "Bearer <token>"
*/

/**
@api {get} /me Request user info
@apiName GetUser
@apiGroup User
@apiUse AuthHeader
@apiUse UserObjectSuccess
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
@apiUse AuthHeader
*/
func (e *Endpoints) getUsers(w http.ResponseWriter, req *http.Request) {
	// TODO: Write docs and tests
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

// Expects a JSON body with the user object with a new password
// returns a 200 status on success
// Deprecated: use UpdateUser instead.
func (e *Endpoints) changePassword(w http.ResponseWriter, r *http.Request) {
	// Decode user struct and check if anything is invalid.
	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil || len(u.Password) == 0 {
		http.Error(w, "Please yeet us a valid body", http.StatusBadRequest)
		return
	}

	password := u.Password

	user := r.Context().Value(contextKeyUser).(db.User)

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


/**
@api {put} /me Update user info
@apiName UpdateUser
@apiGroup User
@apiUse AuthHeader
@apiUse UserObjectParam
@apiUse UserObjectSuccess
@apiError 400 If an invalid body is provided
@apiError 401 If a non-admin changes another user or tries to make themselves admin or blocked
@apiError 422 If the provided password is deemed to weak
*/
func (e *Endpoints) updateUser(w http.ResponseWriter, r *http.Request) {
	// PUT /me (but maybe we need to change it to an actual parameterized route as to make changing username possible
	// and make it a bit more logical for admins as /me doesn't necessarily refer to yourself for them

	u := r.Context().Value(contextKeyUser).(*db.User)
	if u == nil {
		http.Error(w, "This shouldn't happen :tm:", http.StatusInternalServerError)
		return
	}

	var body db.User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Please yeet us a valid body", http.StatusBadRequest)
		return
	}

	if u.Username != body.Username && u.Role != db.AdminRoleID {
		http.Error(w, "Please only edit yourself", http.StatusUnauthorized)
		return
	}

	if body.Blocked == true && u.Role != db.AdminRoleID {
		http.Error(w, "You can't block yourself", http.StatusUnauthorized)
		return
	}

	if body.Role == db.AdminRoleID && u.Role != db.AdminRoleID {
		http.Error(w, "Nice try ;)", http.StatusUnauthorized)
		return
	}

	// If new password provided hash and check
	if len(body.Password) > 0 {
		userinput := []string{
			body.Username,
			body.Email,
		}
		if !passwords.VerifyPassword(body.Password, userinput) {
			http.Error(w, "Please pick a stronger password", http.StatusUnprocessableEntity)
			return
		}

		passwordhash, err := hash.HashPassword(body.Password)
		if err != nil {
			http.Error(w, "Hashing password failed", http.StatusInternalServerError)
			return
		}

		body.Password = passwordhash
	}

	if len(body.Email) > 0 && body.Email != u.Email {
		// TODO: Send new confirmation email
	}


	err := e.conn.UpdateUser(body)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	body.Password = ""

	bytes, err := json.Marshal(body)

	_, err = w.Write(bytes)
	if err != nil {
		log.Error("Couldn't write to client")
	}
}


