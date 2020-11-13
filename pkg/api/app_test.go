package api

import (
	"encoding/json"
	"fmt"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddApplication(t *testing.T) {
	app := models.Application{
		Name:              "app",
		AllowRegistration: false,
	}

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/application", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+tp.LoginToken, token)

		var resp models.Application
		err := json.NewDecoder(r.Body).Decode(&resp)
		assert.NoError(t, err)

		assert.Equal(t, app, resp)

		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	err := AddApplication(ts.URL, &tp, &app)
	assert.NoError(t, err)
}

func TestRemoveApplication(t *testing.T) {
	app := "app"

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/application/"+app, r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+tp.LoginToken, token)
	}))
	defer ts.Close()

	err := RemoveApplication(ts.URL, &tp, app)
	assert.NoError(t, err)
}

func TestGetAccess(t *testing.T) {
	user := "user"
	app := "app"
	url := fmt.Sprintf("/application/%s/%s", app, user)

	access := models.AccessStatus{
		ApplicationName: app,
		Username:        user,
		AllowedAccess:   true,
		Role:            models.RoleAdmin,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, url, r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		err := json.NewEncoder(w).Encode(&access)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	res, err := GetAccess(ts.URL, app, user)
	assert.NoError(t, err)

	assert.Equal(t, access, res)
}

func TestSetAccess(t *testing.T) {
	user := "user"
	app := "app"
	url := fmt.Sprintf("/application/%s/%s", app, user)

	access := models.AccessStatus{
		ApplicationName: app,
		Username:        user,
		AllowedAccess:   true,
		Role:            models.RoleAdmin,
	}

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, url, r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)

		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+tp.LoginToken, token)

		var resp models.AccessStatus
		err := json.NewDecoder(r.Body).Decode(&resp)
		assert.NoError(t, err)

		assert.Equal(t, access, resp)
	}))
	defer ts.Close()

	err := SetAccess(ts.URL, &tp, access)
	assert.NoError(t, err)
}

func TestAddUserToApplication(t *testing.T) {
	user := "user"
	app := "app"
	url := fmt.Sprintf("/application/%s/%s", app, user)

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, url, r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+tp.LoginToken, token)
	}))
	defer ts.Close()

	err := AddUserToApplication(ts.URL, &tp, user, app)
	assert.NoError(t, err)
}

func TestRemoveUserFromApplication(t *testing.T) {
	user := "user"
	app := "app"
	url := fmt.Sprintf("/application/%s/%s", app, user)

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, url, r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+tp.LoginToken, token)
	}))
	defer ts.Close()

	err := RemoveUserFromApplication(ts.URL, &tp, user, app)
	assert.NoError(t, err)
}
