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

func TestAddGroup(t *testing.T) {
	group := models.Group{
		Name:              "group",
		AllowRegistration: false,
	}

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/group", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+tp.LoginToken, token)

		var resp models.Group
		err := json.NewDecoder(r.Body).Decode(&resp)
		assert.NoError(t, err)

		assert.Equal(t, group, resp)

		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	err := AddGroup(ts.URL, &tp, &group)
	assert.NoError(t, err)
}

func TestRemoveGroup(t *testing.T) {
	group := "group"

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/group/"+group, r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+tp.LoginToken, token)
	}))
	defer ts.Close()

	err := RemoveGroup(ts.URL, &tp, group)
	assert.NoError(t, err)
}

func TestGetAccess(t *testing.T) {
	user := "user"
	group := "group"
	url := fmt.Sprintf("/group/%s/%s", group, user)

	access := models.AccessStatus{
		GroupName:     group,
		Username:      user,
		AllowedAccess: true,
		Role:          models.RoleAdmin,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, url, r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		err := json.NewEncoder(w).Encode(&access)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	res, err := GetAccess(ts.URL, group, user)
	assert.NoError(t, err)

	assert.Equal(t, access, res)
}

func TestSetAccess(t *testing.T) {
	user := "user"
	group := "group"
	url := fmt.Sprintf("/group/%s/%s", group, user)

	access := models.AccessStatus{
		GroupName:     group,
		Username:      user,
		AllowedAccess: true,
		Role:          models.RoleAdmin,
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

func TestAddUserToGroup(t *testing.T) {
	user := "user"
	group := "group"
	url := fmt.Sprintf("/group/%s/%s", group, user)

	tp := jwt.TokenPair{
		LoginToken:   "login",
		RefreshToken: "refresh",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, url, r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		token := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+tp.LoginToken, token)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	err := AddUserToGroup(ts.URL, &tp, user, group)
	assert.NoError(t, err)
}

func TestRemoveUserFromGroup(t *testing.T) {
	user := "user"
	group := "group"
	url := fmt.Sprintf("/group/%s/%s", group, user)

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

	err := RemoveUserFromGroup(ts.URL, &tp, user, group)
	assert.NoError(t, err)
}
