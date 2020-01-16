package web

import (
	"aurum/config"
	"aurum/db"
	"aurum/jwt"
	"context"
	"encoding/json"
	_ "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.getMe(w, r.WithContext(ctx))

	// Collect response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ru := db.User{}

	err = json.NewDecoder(resp.Body).Decode(&ru)
	assert.NoError(t, err)
	assert.Equal(t, u, ru)
}
