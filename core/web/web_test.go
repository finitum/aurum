package web

import (
	"aurum/config"
	"aurum/db"
	"aurum/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticateRequest(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.GetDefault()
	endpoints := Endpoints{conn, cfg}

	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.UserRoleID,
		Blocked:  false,
	}

	tkn, err := jwt.GenerateJWT(&u, false, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	claims, err := endpoints.authenticateRequest(w, r)
	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
	assert.Equal(t, u.Username, claims.Username)
}

func TestAuthenticateRequestRefreshToken(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.GetDefault()
	endpoints := Endpoints{conn, cfg}

	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.UserRoleID,
		Blocked:  false,
	}

	tkn, err := jwt.GenerateJWT(&u, true, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	claims, err := endpoints.authenticateRequest(w, r)
	assert.Error(t, err)
	assert.Empty(t, claims)
}

func TestAuthenticateRequestNoToken(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.GetDefault()
	endpoints := Endpoints{conn, cfg}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer ")

	claims, err := endpoints.authenticateRequest(w, r)
	assert.Error(t, err)
	assert.Empty(t, claims)
}

func TestAuthenticateRequestNoHeader(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.GetDefault()
	endpoints := Endpoints{conn, cfg}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Del("Authorization") // Make sure header doesn't exist

	claims, err := endpoints.authenticateRequest(w, r)
	assert.Error(t, err)
	assert.Empty(t, claims)
}

func TestAuthenticateRequestInvalidToken(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.GetDefault()
	endpoints := Endpoints{conn, cfg}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer "+"invalidtkn=")

	claims, err := endpoints.authenticateRequest(w, r)
	assert.Error(t, err)
	assert.Empty(t, claims)
}
