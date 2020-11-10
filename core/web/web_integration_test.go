package web

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/finitum/aurum/core/config"
	"github.com/finitum/aurum/core/db"
	"github.com/finitum/aurum/internal/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/test-go/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TODO: Same as above but for admin user and admin capabilities
func TestSignupLoginFlowIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	config := config.EphemeralConfig()
	repos := db.InitDB(db.InMemory)
	endpoints := Endpoints{
		Repos:  repos,
		Config: config,
	}

	u := models.User{
		Username: "Test",
		Email:    "Tester@test.com",
		Password: "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
		Blocked:  false,
	}

	body, err := json.Marshal(u)
	assert.NoError(t, err)

	// Signup
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))

	endpoints.Signup(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	body, err = json.Marshal(u)
	assert.NoError(t, err)

	// We are first user so we should be admin from here on out
	u.Role = models.AdminRoleID

	// Login
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))

	endpoints.Login(w, r)

	res = w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var token jwt.TokenPair
	err = json.NewDecoder(res.Body).Decode(&token)
	assert.NoError(t, err)

	// Authenticated route
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/user", nil)

	ctx := r.Context()
	ctx = context.WithValue(ctx, "aurum web context key user", &u)

	endpoints.GetMe(w, r.WithContext(ctx))
	res = w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var newuser models.User
	err = json.NewDecoder(res.Body).Decode(&newuser)
	assert.NoError(t, err)
	assert.Equal(t, u.Username, newuser.Username)
	assert.Equal(t, u.Blocked, newuser.Blocked)
	assert.Equal(t, models.AdminRoleID, newuser.Role)
	assert.Equal(t, u.Email, newuser.Email)
}

func VerifyLogin(assert *assert.Assertions, client *http.Client, u models.User) jwt.TokenPair {
	body, err := json.Marshal(u)
	assert.NoError(err)

	req, err := http.NewRequest("POST", "http://localhost:40152/login", bytes.NewBuffer(body))
	assert.NoError(err)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var tp jwt.TokenPair
	err = json.NewDecoder(resp.Body).Decode(&tp)

	assert.NotEmpty(tp.LoginToken)
	assert.NotEmpty(tp.RefreshToken)

	return tp
}

func VerifySignupLogin(assert *assert.Assertions, client *http.Client, u models.User) jwt.TokenPair {

	body, err := json.Marshal(u)
	assert.NoError(err)

	// Signup
	req, err := http.NewRequest("POST", "http://localhost:40152/signup", bytes.NewBuffer(body))
	assert.NoError(err)

	resp, err := client.Do(req)
	assert.NoError(err)

	assert.Equal(http.StatusCreated, resp.StatusCode)

	// Login
	return VerifyLogin(assert, client, u)
}

func VerifyGetUser(assert *assert.Assertions, client *http.Client, tp jwt.TokenPair, u models.User) {

	// get user
	req, err := http.NewRequest("GET", "http://localhost:40152/user", nil)
	assert.NoError(err)
	req.Header.Add("Authorization", "Bearer "+tp.LoginToken)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var newuser models.User
	err = json.NewDecoder(resp.Body).Decode(&newuser)

	assert.Equal(u.Username, newuser.Username)
	assert.Equal(u.Blocked, newuser.Blocked)
	assert.Equal(u.Role, newuser.Role)
	assert.Equal(u.Email, newuser.Email)

}

func VerifyRefresh(assert *assert.Assertions, client *http.Client, tp jwt.TokenPair, u models.User, cfg *config.Config) {
	oldClaims, err := jwt.VerifyJWT(tp.LoginToken, cfg.PublicKey)
	assert.NoError(err)

	body, err := json.Marshal(tp)
	assert.NoError(err)

	// Wait so that the refresh token definitely should have a higher iat
	time.Sleep(2 * time.Second)

	// Refresh
	req, err := http.NewRequest("POST", "http://localhost:40152/refresh", bytes.NewBuffer(body))
	assert.NoError(err)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var rtp jwt.TokenPair
	err = json.NewDecoder(resp.Body).Decode(&rtp)
	assert.Empty(rtp.RefreshToken)

	newClaims, err := jwt.VerifyJWT(rtp.LoginToken, cfg.PublicKey)
	assert.NoError(err)

	assert.True(oldClaims.IssuedAt < newClaims.IssuedAt)

	tp.LoginToken = rtp.LoginToken
	VerifyGetUser(assert, client, tp, u)
}

func VerifyUpdateUserPasswordEmail(assert *assert.Assertions, client *http.Client, tp jwt.TokenPair, u models.User) {
	newuser := models.User{
		Username: u.Username,
		Password: "9054fbe0b622c638224d50d20824d2ff6782e308",
		Email:    "yeet42@finitum.dev",
	}

	body, err := json.Marshal(newuser)
	assert.NoError(err)

	req, err := http.NewRequest("PUT", "http://localhost:40152/user", bytes.NewBuffer(body))
	assert.NoError(err)
	req.Header.Add("Authorization", "Bearer "+tp.LoginToken)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var respuser models.User
	err = json.NewDecoder(resp.Body).Decode(&respuser)

	assert.Equal(u.Username, respuser.Username)
	assert.Equal(newuser.Email, respuser.Email)
	assert.Empty(respuser.Password)

	u.Password = newuser.Password
	u.Email = respuser.Email

	VerifyLogin(assert, client, u)
	VerifyGetUser(assert, client, tp, u)
}

// Warning: this blocks the userToBlock
func VerifyBlockUser(assert *assert.Assertions, client *http.Client, tp jwt.TokenPair, userToBlock models.User) {
	blockedUser := models.User{
		Username: userToBlock.Username,
		Blocked:  true,
	}

	body, err := json.Marshal(blockedUser)
	assert.NoError(err)

	req, err := http.NewRequest("PUT", "http://localhost:40152/user", bytes.NewBuffer(body))
	assert.NoError(err)
	req.Header.Add("Authorization", "Bearer "+tp.LoginToken)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var respuser models.User
	err = json.NewDecoder(resp.Body).Decode(&respuser)

	assert.True(respuser.Blocked)

	// Check the user we just blocked can't login
	blockedUserLogin := models.User{
		Username: userToBlock.Username,
		Password: userToBlock.Password,
	}

	body, err = json.Marshal(blockedUserLogin)
	assert.NoError(err)

	req, err = http.NewRequest("POST", "http://localhost:40152/login", bytes.NewBuffer(body))
	assert.NoError(err)

	resp, err = client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func VerifyOptionsHeaders(assert *assert.Assertions, client *http.Client) {
	req, err := http.NewRequest("OPTIONS", "http://localhost:40152/user", nil)
	assert.NoError(err)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	assert.Equal("*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal("GET, POST, OPTIONS, PUT, DELETE", resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal("Origin, Content-Type, Authorization", resp.Header.Get("Access-Control-Allow-Headers"))
}

func VerifyGetUsers(assert *assert.Assertions, client *http.Client, tp jwt.TokenPair, users []models.User) {
	req, err := http.NewRequest("GET", "http://localhost:40152/users", nil)
	assert.NoError(err)

	rg := Range{0, 100}
	req.URL.RawQuery = rg.toQueryParameters()

	req.Header.Add("Authorization", "Bearer "+tp.LoginToken)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var ub []models.User
	err = json.NewDecoder(resp.Body).Decode(&ub)
	assert.NoError(err)

	assert.Equal(users[0].Username, ub[0].Username)
	assert.Equal(models.AdminRoleID, ub[0].Role)
	assert.Empty(ub[0].Password)
	assert.Equal(users[1].Username, ub[1].Username)
}

func TestSystemIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	assert := assert.New(t)

	cfg := config.EphemeralConfig()
	cfg.WebAddr = "0.0.0.0:40152"

	database := db.InitDB(db.InMemory)

	// Startup the server
	go StartServer(cfg, database)

	// totally not flaky or something
	// Wait for the server  to start up
	time.Sleep(3 * time.Second)

	admin := models.User{
		Username: "TestAdmin",
		Email:    "Tester@test.com",
		Password: "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
		Blocked:  false,
	}

	normal := models.User{
		Username: "TestNormal",
		Email:    "Tester@test.com",
		Password: "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
		Blocked:  false,
	}

	client := &http.Client{}

	// Now run all the endpoint verifications

	VerifyOptionsHeaders(assert, client)

	tpadmin := VerifySignupLogin(assert, client, admin)
	tpnormal := VerifySignupLogin(assert, client, normal)

	admin.Role = models.AdminRoleID
	VerifyGetUser(assert, client, tpadmin, admin)
	VerifyGetUser(assert, client, tpnormal, normal)

	VerifyRefresh(assert, client, tpadmin, admin, cfg)
	VerifyRefresh(assert, client, tpnormal, normal, cfg)

	VerifyUpdateUserPasswordEmail(assert, client, tpnormal, normal)

	users := []models.User{
		admin,
		normal,
	}

	VerifyGetUsers(assert, client, tpadmin, users)

	// after this "normal" is blocked
	VerifyBlockUser(assert, client, tpadmin, normal)
}
