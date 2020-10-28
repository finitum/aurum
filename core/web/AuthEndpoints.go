package web

import (
	"aurum/db"
	"aurum/hash"
	"aurum/jwt"
	"aurum/passwords"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

/**
@api {post} /signup Register
@apiDescription Creates a new account
@apiName Signup
@apiGroup Authentication
@apiParam {String} username The username of the user
@apiParam {String} password The password of the user
@apiParam {String} email The E-Mail of the user
@apiParamExample {json} Request Example:
	{
		"username": "victor",
		"password": "hunter2",
		"email": "victor@example.com"
	}
@apiSuccessExample {String} Success Response:
	HTTP/1.1 201 Created

@apiError 400 If an invalid body is provided
@apiError 422 If an insufficiently secure password is provided
*/
func (e *Endpoints) Signup(w http.ResponseWriter, r *http.Request) {

	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		if err != nil {
			log.Error(err)
		}
		http.Error(w, "Please yeet us a valid json body", http.StatusBadRequest)
		return
	}

	if u.Username == "" || u.Password == "" || u.Email == "" || u.Role != db.UserRoleID {
		http.Error(w, "Please yeet us a valid json body", http.StatusBadRequest)
		return
	}

	if !passwords.VerifyPassword(u.Password, []string{u.Username, u.Email}) {
		http.Error(w, "Password not acceptable", http.StatusUnprocessableEntity)
		return
	}

	// If the signed up user is the only user, make this user admin
	number, err := e.Repos.CountUsers()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if number == 0 {
		log.Infof("A user with username \"%s\" has signed up. This is the first user and will get admin privileges.", u.Username)
		u.Role = db.AdminRoleID
	}

	// Actually add the user to the db
	err = e.Repos.CreateUser(u)
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

/**
@api {post} /refresh Refresh Token
@apiDescription Refreshes your Login token by using your refresh token
@apiName Refresh
@apiGroup Authentication
@apiParam {String} refresh_token The refresh token to use.
@apiParamExample {json} Request Example:
	{
		"refresh_token": "<JWT Token here>"
	}
@apiSuccess {String} login_token A renewed Login token
@apiSuccessExample {json} Success Response:
	{
		"login_token": "<JWT Token here>"
	}

@apiError 400 If an invalid body or token is provided
@apiError 404 If the user does not exist (anymore)
*/
func (e *Endpoints) Refresh(w http.ResponseWriter, r *http.Request) {

	var t jwt.TokenPair
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Please yeet us a valid json body", http.StatusBadRequest)
		return
	}

	c, err := jwt.VerifyJWT(t.RefreshToken, e.Config)
	if err != nil || !c.Refresh {
		http.Error(w, "Please specify a valid refresh token", http.StatusBadRequest)
		return
	}

	user, err := e.Repos.GetUserByName(c.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	token, err := jwt.GenerateJWT(&user, false, e.Config)
	if err != nil {
		http.Error(w, "Couldn't generate Login token", http.StatusInternalServerError)
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

/**
@api {post} /Login Login
@apiDescription Logs a user in returning a tokenpair
@apiName Login
@apiGroup Authentication
@apiParam {String} username The user's username
@apiParam {String} password The user's password
@apiParamExample {json} Request Example:
	{
		"username": "victor",
		"password": "hunter2"
	}
@apiSuccess {String} login_token The user's Login token
@apiSuccess {String} refresh_token The user's refresh token
@apiSuccessExample {json} Success Response:
	{
		"login_token": "<JWT Token here>"
		"refresh_token": "<JWT Token here>"
	}
@apiError 400 If an invalid body is provided.
@apiError 401 If the user does not exist or the password is wrong
*/
func (e *Endpoints) Login(w http.ResponseWriter, r *http.Request) {

	// Decode user struct and check if anything is invalid.
	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil || len(u.Username) == 0 || len(u.Password) == 0 {
		if err != nil {
			log.Error(err)
		}
		http.Error(w, "Please yeet us a valid body", http.StatusBadRequest)
		return
	}

	user, err := e.Repos.GetUserByName(u.Username)
	if err != nil || user.Blocked {
		if err != nil {
			log.Error(err)
		}
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}

	// Invalid password
	if !hash.CheckPasswordHash(u.Password, user.Password) {
		http.Error(w, "User is not authorized", http.StatusUnauthorized)
		return
	}

	// Generate a JWT pair (Login + refresh)
	token, err := jwt.GenerateJWTPair(&user, e.Config)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error in JWT generation", http.StatusInternalServerError)
		return
	}

	// Convert the tokens into json bytes
	bytes, err := json.Marshal(token)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error in serializing JWT", http.StatusInternalServerError)
		return
	}

	// push them to the client
	if _, err = w.Write(bytes); err != nil {
		log.Error("Error writing response to client")
	}

	log.Infof("A user with username %v has logged in.", u.Username)

	return
}

type publicKeyResponse struct {
	PublicKey string `json:"public_key"`
}

/**
@api {get} /pk PublicKey
@apiDescription Returns the server's publickey
@apiName PublicKey
@apiSuccess {String} public_key the base64 encoded public key
@apiSuccessExample {json} Success Response:
	{
		"public_key": "<Public key here>"
	}
*/
func (e *Endpoints) PublicKey(w http.ResponseWriter, r *http.Request) {
	pem, err := e.Config.PublicKey.ToPem()
	if err != nil {
		http.Error(w, "Error in getting PEM", http.StatusInternalServerError)
		return
	}

	pk := publicKeyResponse{
		PublicKey: pem,
	}

	bytes, err := json.Marshal(pk)

	if err != nil {
		http.Error(w, "Error in serializing PK", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err = w.Write(bytes); err != nil {
		log.Error("Error writing response to client")
	}

	return
}
