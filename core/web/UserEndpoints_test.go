package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/finitum/aurum/core/config"
	"github.com/finitum/aurum/internal/jwt"
	"github.com/finitum/aurum/pkg/hash"
	"github.com/finitum/aurum/pkg/models"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMe(t *testing.T) {
	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	conn := SQLConnectionMock{}

	conn.On("GetUserByName", u.Username).Return(u, nil)

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.GetMe(w, r.WithContext(ctx))

	// Collect response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ru := models.User{}

	err = json.NewDecoder(resp.Body).Decode(&ru)
	assert.NoError(t, err)
	assert.Equal(t, u, ru)
}

func TestChangePassword(t *testing.T) {
	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Password = "cf7c1ee3fc9d90ec191ae576eb851022a7414741"

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}

	conn.On("UpdateUser", mock.Anything).Run(func(args mock.Arguments) {
		user := args.Get(0).(models.User)
		assert.Equal(t, u.Username, user.Username)
		assert.Equal(t, u.Email, user.Email)
		assert.Equal(t, u.Role, user.Role)
		assert.Equal(t, u.Blocked, user.Blocked)
		assert.True(t, hash.CheckPasswordHash(uModified.Password, user.Password))
	}).Return(nil)

	cfg := config.EphemeralConfig()

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.UpdateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var ub models.User
	_ = json.Unmarshal(w.Body.Bytes(), &ub)
	assert.Equal(t, u, ub)

	conn.AssertExpectations(t)
}

func TestBlockSelf(t *testing.T) {
	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Blocked = true

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.UpdateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	conn.AssertExpectations(t)
}

func TestAdminSelf(t *testing.T) {
	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Role = models.AdminRoleID

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.UpdateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	conn.AssertExpectations(t)
}

func TestChangePasswordOtherUserAsAdmin(t *testing.T) {
	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.AdminRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Username = "notme"
	uModified.Password = "cf7c1ee3fc9d90ec191ae576eb851022a7414741"

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}

	conn.On("UpdateUser", mock.Anything).Run(func(args mock.Arguments) {
		user := args.Get(0).(models.User)
		assert.Equal(t, uModified.Username, user.Username)
		assert.Equal(t, u.Email, user.Email)
		assert.Equal(t, u.Role, user.Role)
		assert.Equal(t, u.Blocked, user.Blocked)
		assert.True(t, hash.CheckPasswordHash(uModified.Password, user.Password))
	}).Return(nil)

	cfg := config.EphemeralConfig()

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.UpdateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var ub models.User
	_ = json.Unmarshal(w.Body.Bytes(), &ub)
	uModified.Password = ""
	assert.Equal(t, uModified, ub)

	conn.AssertExpectations(t)
}

func TestChangeWrongUsername(t *testing.T) {
	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	uModified := u
	uModified.Username = "notme"
	uModified.Password = "cf7c1ee3fc9d90ec191ae576eb851022a7414741"

	body, err := json.Marshal(uModified)

	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewBuffer(body))
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.UpdateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	conn.AssertExpectations(t)
}

func TestChangeUserNoBody(t *testing.T) {
	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	endpoints.UpdateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	conn.AssertExpectations(t)
}

type FakeWriter struct{}

func (f FakeWriter) Header() http.Header {
	panic("implement me")
}

func (f FakeWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func (f FakeWriter) WriteHeader(_ int) {
	panic("implement me")
}

func TestChangeUserBadWriter(t *testing.T) {
	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	w := FakeWriter{}

	// install new log hook
	hook := test.NewGlobal()

	endpoints.GetMe(w, r.WithContext(ctx))

	// assert that there was an error log
	assert.Equal(t, len(hook.Entries), 1)
	assert.Equal(t, hook.Entries[0].Level, log.ErrorLevel)
	conn.AssertExpectations(t)

}

func TestGetUsersNotAdmin(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{&conn, cfg}

	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	w := httptest.NewRecorder()

	endpoints.GetUsers(w, r.WithContext(ctx))

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)

	conn.AssertExpectations(t)

}

func TestGetUsersInvalidRange(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{&conn, cfg}

	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.AdminRoleID,
		Blocked:  false,
	}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	rg := Range{100, 0}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)
	r.URL.RawQuery = rg.toQueryParameters()

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	w := httptest.NewRecorder()

	endpoints.GetUsers(w, r.WithContext(ctx))

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	conn.AssertExpectations(t)
}

func TestGetUsersDbError(t *testing.T) {
	conn := SQLConnectionMock{}

	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.AdminRoleID,
		Blocked:  false,
	}

	cfg := config.EphemeralConfig()

	conn.On("GetUserByName", u.Username).Return(u, nil)
	conn.On("GetUsers", mock.Anything, mock.Anything).Return([]models.User{}, errors.New("simulated DB Error"))

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	rg := Range{0, 100}
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)
	r.URL.RawQuery = rg.toQueryParameters()

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	w := httptest.NewRecorder()

	endpoints.GetUsers(w, r.WithContext(ctx))

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
}

func TestGetUsers(t *testing.T) {
	conn := SQLConnectionMock{}

	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.AdminRoleID,
		Blocked:  false,
	}

	cfg := config.EphemeralConfig()

	conn.On("GetUsers", 0, 100).Return([]models.User{u}, nil)

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	rg := Range{0, 100}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)
	r.URL.RawQuery = rg.toQueryParameters()

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	w := httptest.NewRecorder()

	endpoints.GetUsers(w, r.WithContext(ctx))

	resp := w.Result()

	assert.Equal(t, resp.StatusCode, http.StatusOK)

	var ub []models.User
	_ = json.Unmarshal(w.Body.Bytes(), &ub)
	assert.Equal(t, []models.User{u}, ub)

	conn.AssertExpectations(t)
}

func TestGetUsersBadWriter(t *testing.T) {
	conn := SQLConnectionMock{}

	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.AdminRoleID,
		Blocked:  false,
	}

	cfg := config.EphemeralConfig()

	conn.On("GetUsers", mock.Anything, mock.Anything).Return([]models.User{u}, nil)

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	rg := Range{0, 100}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)
	r.URL.RawQuery = rg.toQueryParameters()

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, &u)

	w := FakeWriter{}

	// install new log hook
	hook := test.NewGlobal()

	endpoints.GetUsers(w, r.WithContext(ctx))

	// assert that there was an error log
	assert.Equal(t, len(hook.Entries), 1)
	assert.Equal(t, hook.Entries[0].Level, log.ErrorLevel)

	conn.AssertExpectations(t)
}

func TestUpdateUserNilUser(t *testing.T) {
	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	conn := SQLConnectionMock{}
	cfg := config.EphemeralConfig()

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	ctx := r.Context()
	ctx = context.WithValue(ctx, contextKeyUser, nil)

	endpoints.UpdateUser(w, r.WithContext(ctx))

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	conn.AssertExpectations(t)

}
