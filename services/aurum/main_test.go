package main

import (
	"bytes"
	"encoding/json"
	"github.com/finitum/aurum/internal/aurum"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"github.com/finitum/aurum/pkg/models"
)

const url = "http://localhost:8042"

func VerifyLogin(assert *assert.Assertions, client *http.Client, u models.User) jwt.TokenPair {
	body, err := json.Marshal(u)
	assert.NoError(err)

	req, err := http.NewRequest("POST", url+"/login", bytes.NewBuffer(body))
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
	req, err := http.NewRequest("POST", url+"/signup", bytes.NewBuffer(body))
	assert.NoError(err)

	resp, err := client.Do(req)
	assert.NoError(err)

	assert.Equal(http.StatusCreated, resp.StatusCode)

	time.Sleep(time.Second)

	// Login
	return VerifyLogin(assert, client, u)
}

func VerifyGetUser(assert *assert.Assertions, client *http.Client, tp jwt.TokenPair, u models.User) {

	// get user
	req, err := http.NewRequest("GET", url+"/user", nil)
	assert.NoError(err)
	req.Header.Add("Authorization", "Bearer "+tp.LoginToken)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var newuser models.User
	err = json.NewDecoder(resp.Body).Decode(&newuser)

	assert.Equal(u.Username, newuser.Username)
	assert.Equal(u.Email, newuser.Email)

}

func VerifyRefresh(assert *assert.Assertions, client *http.Client, tp jwt.TokenPair, u models.User, pk ecc.PublicKey) {
	oldClaims, err := jwt.VerifyJWT(tp.LoginToken, pk)
	assert.NoError(err)

	body, err := json.Marshal(tp)
	assert.NoError(err)

	// Wait so that the refresh token definitely should have a higher iat
	time.Sleep(2 * time.Second)

	// Refresh
	req, err := http.NewRequest("POST", url+"/refresh", bytes.NewBuffer(body))
	assert.NoError(err)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var rtp jwt.TokenPair
	err = json.NewDecoder(resp.Body).Decode(&rtp)
	assert.Empty(rtp.RefreshToken)

	newClaims, err := jwt.VerifyJWT(rtp.LoginToken, pk)
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

	req, err := http.NewRequest("POST", url+"/user", bytes.NewBuffer(body))
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
	u.Email = newuser.Email

	time.Sleep(2 * time.Second)

	VerifyLogin(assert, client, u)
	VerifyGetUser(assert, client, tp, u)
}

func VerifyGetGroupsForUser(assert *assert.Assertions, client *http.Client, tp jwt.TokenPair, u models.User, expected models.Role) {
	req, err := http.NewRequest(http.MethodGet, url+"/user/"+u.Username+"/groups", nil)
	assert.NoError(err)
	req.Header.Add("Authorization", "Bearer "+tp.LoginToken)

	resp, err := client.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	var groups []models.GroupWithRole
	err = json.NewDecoder(resp.Body).Decode(&groups)
	assert.NoError(err)

	aurums := 0
	for _, group := range groups {
		if group.Name == aurum.AurumName {
			assert.Equal(expected, group.Role)
			aurums += 1
		}
	}
	assert.Equal(1, aurums)
}

func TestSystemIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	assert := assert.New(t)

	assert.NoError(os.Setenv("NO_KEY_WRITE", "true"))
	assert.NoError(os.Setenv("WEB_ADDRESS", strings.TrimPrefix(url, "http://")))

	// Startup the server
	go main()

	// totally not flaky or something
	// Wait for the server  to start up
	time.Sleep(5 * time.Second)

	userOne := models.User{
		Username: "UserOne",
		Email:    "Tester@test.com",
		Password: "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
	}

	userTwo := models.User{
		Username: "UserTwo",
		Email:    "Tester@test.com",
		Password: "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
	}

	client := &http.Client{}

	// Now run all the endpoint verifications

	resp, err := http.Get(url + "/pk")
	assert.NoError(err)

	var r models.PublicKeyResponse

	err = json.NewDecoder(resp.Body).Decode(&r)
	assert.NoError(err)

	pk, err := ecc.FromPem([]byte(r.PublicKey))
	assert.NoError(err)

	pub := pk.(ecc.PublicKey)

	tpUserOne := VerifySignupLogin(assert, client, userOne)
	tpUserTwo := VerifySignupLogin(assert, client, userTwo)

	time.Sleep(time.Second)

	VerifyGetUser(assert, client, tpUserOne, userOne)
	VerifyGetUser(assert, client, tpUserTwo, userTwo)

	VerifyRefresh(assert, client, tpUserOne, userOne, pub)
	VerifyRefresh(assert, client, tpUserTwo, userTwo, pub)

	VerifyUpdateUserPasswordEmail(assert, client, tpUserOne, userOne)
	VerifyUpdateUserPasswordEmail(assert, client, tpUserTwo, userTwo)

	VerifyGetGroupsForUser(assert, client, tpUserOne, userOne, models.RoleUser)
	VerifyGetGroupsForUser(assert, client, tpUserTwo, userTwo, models.RoleUser)
}
