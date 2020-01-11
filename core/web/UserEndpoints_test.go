package web

import (
	"aurum/config"
	"aurum/db"
	"aurum/jwt"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMe(t *testing.T) {
	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.UserRoleID,
		Blocked:  false,
	}

	conn := SQLConnectionMock{}

	conn.On("GetUserByName", u.Username).Return(u, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	endpoints.getMe(w, r)

	// Collect response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ru := db.User{}
	err = json.NewDecoder(resp.Body).Decode(&ru)
	assert.NoError(t, err)
	assert.Equal(t, u, ru)
}

func TestGetMeInvalidUser(t *testing.T) {
	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.UserRoleID,
		Blocked:  false,
	}

	conn := SQLConnectionMock{}

	conn.On("GetUserByName", u.Username).Return(db.User{}, errors.New("no user found"))

	cfg := config.GetDefault()
	endpoints := Endpoints{conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	endpoints.getMe(w, r)

	// Collect response
	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(t, string(body), "User not found")
}
