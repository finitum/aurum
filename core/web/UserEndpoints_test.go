package web

import (
	"aurum/config"
	"aurum/db"
	"aurum/hash"
	"aurum/jwt"
	"bytes"
	"context"
	"encoding/json"
	_ "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestChangePassword(t *testing.T) {
	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.UserRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Password = "cf7c1ee3fc9d90ec191ae576eb851022a7414741"

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}

	conn.On("UpdateUser", mock.Anything).Run(func(args mock.Arguments) {
		user := args.Get(0).(db.User)
		assert.Equal(t, u.Username, user.Username)
		assert.Equal(t, u.Email, user.Email)
		assert.Equal(t, u.Role, user.Role)
		assert.Equal(t, u.Blocked, user.Blocked)
		assert.True(t, hash.CheckPasswordHash(uModified.Password, user.Password))
	}).Return(nil)

	cfg := config.GetDefault()

	endpoints := Endpoints{conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.updateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var ub db.User
	_ = json.Unmarshal(w.Body.Bytes(), &ub)
	assert.Equal(t, u, ub)

	conn.AssertExpectations(t)
}

func TestBlockSelf(t *testing.T) {
	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.UserRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Blocked = true

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}

	cfg := config.GetDefault()

	endpoints := Endpoints{conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.updateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	conn.AssertExpectations(t)
}

func TestAdminSelf(t *testing.T) {
	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.UserRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Role = db.AdminRoleID

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}

	cfg := config.GetDefault()

	endpoints := Endpoints{conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.updateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	conn.AssertExpectations(t)
}

func TestChangePasswordOtherUserAsAdmin(t *testing.T) {
	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.AdminRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Username = "notme"
	uModified.Password = "cf7c1ee3fc9d90ec191ae576eb851022a7414741"

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}

	conn.On("UpdateUser", mock.Anything).Run(func(args mock.Arguments) {
		user := args.Get(0).(db.User)
		assert.Equal(t, uModified.Username, user.Username)
		assert.Equal(t, u.Email, user.Email)
		assert.Equal(t, u.Role, user.Role)
		assert.Equal(t, u.Blocked, user.Blocked)
		assert.True(t, hash.CheckPasswordHash(uModified.Password, user.Password))
	}).Return(nil)

	cfg := config.GetDefault()

	endpoints := Endpoints{conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.updateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var ub db.User
	_ = json.Unmarshal(w.Body.Bytes(), &ub)
	uModified.Password = ""
	assert.Equal(t, uModified, ub)

	conn.AssertExpectations(t)
}

func TestChangeWrongUsername(t *testing.T) {
	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.UserRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Username = "notme"
	uModified.Password = "cf7c1ee3fc9d90ec191ae576eb851022a7414741"

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}


	cfg := config.GetDefault()

	endpoints := Endpoints{conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.updateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	conn.AssertExpectations(t)
}



func TestChangeUserNoBody(t *testing.T) {
	u := db.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     db.UserRoleID,
		Blocked:  false,
	}

	conn := SQLConnectionMock{}

	cfg := config.GetDefault()

	endpoints := Endpoints{conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.updateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	conn.AssertExpectations(t)
}