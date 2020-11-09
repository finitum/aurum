package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/finitum/aurum/core/config"
	"github.com/finitum/aurum/core/db"
	hash2 "github.com/finitum/aurum/pkg/hash"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type SQLConnectionMock struct {
	mock.Mock
}

func (conn *SQLConnectionMock) GetUserByName(u string) (models.User, error) {
	args := conn.Called(u)
	return args.Get(0).(models.User), args.Error(1)
}

func (conn *SQLConnectionMock) CreateUser(u models.User) error {
	args := conn.Called(u)
	return args.Error(0)
}

func (conn *SQLConnectionMock) CountUsers() (int, error) {
	conn.Called()
	return 1, nil
}

func (conn *SQLConnectionMock) UpdateUser(user models.User) error {
	args := conn.Called(user)
	return args.Error(0)
}

func (conn *SQLConnectionMock) GetUsers(start int, end int) ([]models.User, error) {
	args := conn.Called(start, end)
	return args.Get(0).([]models.User), args.Error(1)
}

func TestSignup(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", models.User{
		Username: "jonathan",
		Password: "7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f",
		Email:    "yeet@yeet.dev",
	}).Return(nil)
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.EphemeralConfig()

	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",
    "password":"7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f",
    "email":"yeet@yeet.dev"
    }`))

	endpoints.Signup(w, r)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	conn.AssertExpectations(t)
}

func TestSignupIncorrectJson(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()

	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{"))

	endpoints.Signup(w, r)

	conn.AssertExpectations(t)
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupNoBody(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	endpoints.Signup(w, r)

	conn.AssertExpectations(t)
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupNoUsername(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "password":"yeet",
    "email":"yeet@yeet.dev"
}`))

	endpoints.Signup(w, r)

	conn.AssertExpectations(t)
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupNoPassword(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "email":"yeet@yeet.dev"
}`))

	endpoints.Signup(w, r)

	conn.AssertExpectations(t)
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupNoEmail(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeet"
}`))

	endpoints.Signup(w, r)

	conn.AssertExpectations(t)
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupAdmin(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeet",
    "email":"yeet@yeet.dev",
    "role": 1
    }`))

	endpoints.Signup(w, r)

	conn.AssertExpectations(t)
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestDbError(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", mock.Anything).Return(errors.New("errortest"))
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"4b93310ed64ce510889be78f32203f9768c4054b9af08489ed90a59465616ef6",
    "email":"yeet@yeet.dev"
}`))

	endpoints.Signup(w, r)

	conn.AssertExpectations(t)
	assert.Equal(t, w.Result().StatusCode, http.StatusInternalServerError)
}

func TestLogin(t *testing.T) {
	conn := SQLConnectionMock{}
	hash, err := hash2.HashPassword("yeetyeet")
	assert.Nil(t, err)

	u := models.User{
		Username: "jonathan",
		Password: hash,
		Email:    "yeet@yeet.dev",
		Role:     models.UserRoleID,
	}

	conn.On("GetUserByName", mock.Anything).Run(func(args mock.Arguments) {
		assert.Equal(t, args.Get(0), "jonathan")
	}).Return(u, nil)

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeetyeet"
    }`))

	endpoints.Login(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	resp := w.Result()

	tp := jwt.TokenPair{}
	err = json.NewDecoder(resp.Body).Decode(&tp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tp.LoginToken)
	assert.NotEmpty(t, tp.RefreshToken)

	// Login token
	claims, err := jwt.VerifyJWT(tp.LoginToken, cfg.PublicKey)
	if assert.NoError(t, err) {
		assert.NotNil(t, claims)
		assert.Equal(t, "jonathan", claims.Username)
		assert.Equal(t, 0, claims.Role)
		assert.Equal(t, false, claims.Refresh)
	}

	// Refresh token
	claims, err = jwt.VerifyJWT(tp.RefreshToken, cfg.PublicKey)
	if assert.NoError(t, err) {
		assert.NotNil(t, claims)
		assert.Equal(t, "jonathan", claims.Username)
		assert.Equal(t, 0, claims.Role)
		assert.Equal(t, true, claims.Refresh)
	}

	conn.AssertExpectations(t)
}

func TestLoginNoBody(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	endpoints.Login(w, r)

	conn.AssertExpectations(t)
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestLoginIncorrectJson(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{"))

	endpoints.Login(w, r)

	conn.AssertExpectations(t)
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestLoginInvalidPassword(t *testing.T) {
	conn := SQLConnectionMock{}
	hash, err := hash2.HashPassword("woolnoo")
	assert.Nil(t, err)

	conn.On("GetUserByName", mock.Anything).Run(func(args mock.Arguments) {
		assert.Equal(t, args.Get(0), "jonathan")
	}).Return(models.User{
		Username: "jonathan",
		Password: hash,
		Email:    "yeet@yeet.dev",
		Role:     models.UserRoleID,
	}, nil)

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeet"
}`))

	endpoints.Login(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusUnauthorized)
	conn.AssertExpectations(t)
}

func TestLoginInvalidUsername(t *testing.T) {
	conn := db.InitDB(db.InMemory)

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  conn,
		Config: cfg,
	}

	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f",
    "email":"test@test.test"
}`))

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathanx",    
    "password":"7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f"
}`))

	endpoints.Signup(w1, r1)
	endpoints.Login(w2, r2)

	assert.Equal(t, http.StatusCreated, w1.Result().StatusCode)
	assert.Equal(t, http.StatusUnauthorized, w2.Result().StatusCode)
}

func TestLoginUsernameMismatch(t *testing.T) {
	conn := db.InitDB(db.InMemory)

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  conn,
		Config: cfg,
	}

	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f",
    "email":"test@test.test"
}`))

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathanx",    
    "password":"4b93310ed64ce510889be78f32203f9768c4054b9af08489ed90a59465616ef6",
    "email":"test@test.test"
}`))

	w3 := httptest.NewRecorder()
	r3 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathanx",    
    "password":"7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f"
}`))

	endpoints.Signup(w1, r1)
	endpoints.Signup(w2, r2)
	endpoints.Login(w3, r3)

	assert.Equal(t, http.StatusCreated, w1.Result().StatusCode)
	assert.Equal(t, http.StatusCreated, w2.Result().StatusCode)
	assert.Equal(t, http.StatusUnauthorized, w3.Result().StatusCode)
}

func TestSignupExists(t *testing.T) {
	conn := db.InitDB(db.InMemory)

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  conn,
		Config: cfg,
	}

	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f",
    "email":"test@test.test"
    }`))

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"4b93310ed64ce510889be78f32203f9768c4054b9af08489ed90a59465616ef6",
    "email":"test@test.test"
    }`))

	endpoints.Signup(w1, r1)
	endpoints.Signup(w2, r2)

	assert.Equal(t, w1.Result().StatusCode, http.StatusCreated)
	assert.Equal(t, w2.Result().StatusCode, http.StatusConflict)
}

func TestRefresh(t *testing.T) {
	u := models.User{
		Username: "victor",
	}

	conn := SQLConnectionMock{}
	conn.On("GetUserByName", "victor").Return(u, nil)

	cfg := config.EphemeralConfig()

	ep := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, true, cfg.SecretKey)
	assert.NoError(t, err)

	b, err := json.Marshal(jwt.TokenPair{
		RefreshToken: tkn,
	})

	// Create httptest recorder and request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(b))

	// Call the refresh endpoint
	ep.Refresh(w, r)

	// Collect the http result
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	tp := jwt.TokenPair{}
	err = json.NewDecoder(res.Body).Decode(&tp)

	// Check result structure
	assert.NoError(t, err)
	assert.NotEqual(t, tkn, tp.RefreshToken)
	assert.NotEmpty(t, tp.LoginToken)
	assert.Empty(t, tp.RefreshToken)

	// Check token and claim validity
	c, err := jwt.VerifyJWT(tp.LoginToken, cfg.PublicKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, c)
	assert.Equal(t, u.Username, c.Username)

	conn.AssertExpectations(t)
}

func TestEndpoints_PublicKey(t *testing.T) {
	conn := SQLConnectionMock{}

	cfg := config.EphemeralConfig()
	endpoints := Endpoints{
		Repos:  &conn,
		Config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/pk", nil)

	endpoints.PublicKey(w, r)

	conn.AssertExpectations(t)

	res := w.Result()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	pk := models.PublicKeyResponse{}
	err := json.NewDecoder(res.Body).Decode(&pk)
	assert.NoError(t, err)

	pem, err := cfg.PublicKey.ToPem()
	assert.NoError(t, err)
	assert.Contains(t, pem, "PUBLIC KEY")
	epk := models.PublicKeyResponse{PublicKey: pem}

	assert.Equal(t, epk, pk)
}
