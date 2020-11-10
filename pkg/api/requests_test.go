package api

import (
	"encoding/json"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/test-go/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPublicKey(t *testing.T) {
	pkrsp := models.PublicKeyResponse{PublicKey: "apublickey"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/pk", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		err := json.NewEncoder(w).Encode(&pkrsp)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	resp, err := GetPublicKey(ts.URL)
	assert.NoError(t, err)

	assert.Equal(t, &pkrsp, resp)
}

func TestSignUp(t *testing.T) {
	u := models.User{
		Username: "user",
		Password: "pass",
		Email:    "email",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/signup", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var recv models.User
		err := json.NewDecoder(r.Body).Decode(&recv)
		assert.NoError(t, err)

		assert.Equal(t, u, recv)

		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	err := SignUp(ts.URL, u)
	assert.NoError(t, err)
}

func TestLogin(t *testing.T) {
	u := models.User{
		Username: "user",
		Password: "pass",
	}

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/login", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var recv models.User
		err := json.NewDecoder(r.Body).Decode(&recv)
		assert.NoError(t, err)

		assert.Equal(t, u, recv)

		err = json.NewEncoder(w).Encode(&tp)
		assert.NoError(t, err)
		}))
	defer ts.Close()

	rtp, err := Login(ts.URL, u)
	assert.NoError(t, err)
	assert.Equal(t, &tp, rtp)
}

func TestRefresh(t *testing.T) {
	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	tp2 := jwt.TokenPair{
		LoginToken:   "login2",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/refresh", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)


		var recv jwt.TokenPair
		err := json.NewDecoder(r.Body).Decode(&recv)
		assert.NoError(t, err)

		assert.Equal(t, tp, recv)

		err = json.NewEncoder(w).Encode(&tp2)
		assert.NoError(t, err)

	}))
	defer ts.Close()

	err := Refresh(ts.URL, &tp)
	assert.NoError(t, err)

	assert.Equal(t, tp2.LoginToken, tp.LoginToken)
	assert.Equal(t, "refresh", tp.RefreshToken)
}

func TestGetUser(t *testing.T) {

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	u := models.User{
		Username: "user",
		Password: "pass",
		Email:    "mail",
		Role:     models.AdminRoleID,
		Blocked:  false,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/user", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)


		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer " + tp.LoginToken, token)

		err := json.NewEncoder(w).Encode(&u)
		assert.NoError(t, err)

	}))
	defer ts.Close()

	user, err := GetUser(ts.URL, &tp)
	assert.NoError(t, err)
	assert.Equal(t, &u, user)
}

func TestUpdateUser(t *testing.T) {

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	u := models.User{
		Username: "user",
		Password: "pass",
		Email:    "mail",
		Role:     models.AdminRoleID,
		Blocked:  false,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/user", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer " + tp.LoginToken, token)

		err := json.NewEncoder(w).Encode(&u)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	user, err := UpdateUser(ts.URL, &tp, &u)
	assert.NoError(t, err)
	assert.Equal(t, &u, user)
}

func TestUpdateUserRefreshNeeded(t *testing.T) {

	const initialLogin = "login"
	const refreshLogin = "login2"

	initialToken := jwt.TokenPair{
		LoginToken:   initialLogin,
		RefreshToken: "refresh",
	}

	refreshToken := jwt.TokenPair{
		LoginToken: refreshLogin,
	}

	u := models.User{
		Username: "user",
		Password: "pass",
		Email:    "mail",
		Role:     models.AdminRoleID,
		Blocked:  false,
	}
	upd := u

	hasRefreshed := false
	hasGot := false
	hasTried := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/refresh" {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.False(t, hasRefreshed)
			assert.False(t, hasGot)
			assert.True(t, hasTried)

			var recv jwt.TokenPair
			err := json.NewDecoder(r.Body).Decode(&recv)
			assert.NoError(t, err)

			assert.Equal(t, initialToken, recv)

			err = json.NewEncoder(w).Encode(&refreshToken)
			assert.NoError(t, err)
			hasRefreshed = true
			return
		} else if r.URL.Path == "/user" {
			assert.Equal(t, http.MethodPost, r.Method)
			token := r.Header.Get("Authorization")
			if token == "Bearer " + initialLogin {
				hasTried = true
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if token == "Bearer " + refreshLogin {
				assert.False(t, hasGot)
				assert.True(t, hasRefreshed)

				var recv models.User
				err := json.NewDecoder(r.Body).Decode(&recv)
				assert.NoError(t, err)
				assert.Equal(t, upd, recv)

				err = json.NewEncoder(w).Encode(&u)
				assert.NoError(t, err)
				hasGot = true
				return
			}
		}
		t.Fail()
	}))
	defer ts.Close()

	user, err := UpdateUser(ts.URL, &initialToken, &upd)
	assert.NoError(t, err)
	assert.Equal(t, &u, user)
}

func TestGetUserRefreshNeeded(t *testing.T) {

	const initialLogin = "login"
	const refreshLogin = "login2"

	initialToken := jwt.TokenPair{
		LoginToken:   initialLogin,
		RefreshToken: "refresh",
	}

	refreshToken := jwt.TokenPair{
		LoginToken: refreshLogin,
	}

	u := models.User{
		Username: "user",
		Password: "pass",
		Email:    "mail",
		Role:     models.AdminRoleID,
		Blocked:  false,
	}

	hasRefreshed := false
	hasGot := false
	hasTried := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/refresh" {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.False(t, hasRefreshed)
			assert.False(t, hasGot)
			assert.True(t, hasTried)

			var recv jwt.TokenPair
			err := json.NewDecoder(r.Body).Decode(&recv)
			assert.NoError(t, err)

			assert.Equal(t, initialToken, recv)

			err = json.NewEncoder(w).Encode(&refreshToken)
			assert.NoError(t, err)
			hasRefreshed = true
			return
		} else if r.URL.Path == "/user" {
			assert.Equal(t, http.MethodGet, r.Method)
			token := r.Header.Get("Authorization")
			if token == "Bearer " + initialLogin {
				hasTried = true
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if token == "Bearer " + refreshLogin {
				assert.False(t, hasGot)
				assert.True(t, hasRefreshed)
				err := json.NewEncoder(w).Encode(&u)
				assert.NoError(t, err)
				hasGot = true
				return
			}
		}
		t.Fail()
	}))
	defer ts.Close()

	user, err := GetUser(ts.URL, &initialToken)
	assert.NoError(t, err)
	assert.Equal(t, &u, user)
}
