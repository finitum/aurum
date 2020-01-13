package web

import (
	"aurum/config"
	"aurum/db"
	hash2 "aurum/hash"
	"aurum/jwt"
	"bytes"
	"encoding/json"
	"errors"
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

func (conn SQLConnectionMock) GetUserByName(u string) (db.User, error) {
	args := conn.Called(u)
	return args.Get(0).(db.User), args.Error(1)
}

func (conn SQLConnectionMock) CreateUser(u db.User) error {
	args := conn.Called(u)
	return args.Error(0)
}

func (conn SQLConnectionMock) CountUsers() (int, error) {
	conn.Called()
	return 1, nil
}

func (conn SQLConnectionMock) UpdateUser(user db.User) error {
	args := conn.Called(user)
	return args.Error(0)
}

func (conn SQLConnectionMock) GetUsers(start int, end int) ([]db.User, error) {
	args := conn.Called(start, end)
	return args.Get(0).([]db.User), args.Error(1)
}

func TestSignup(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", db.User{
		Username: "jonathan",
		Password: "yeetyeet",
		Email:    "yeet@yeet.dev",
	}).Return(nil)
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()

	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",
    "password":"yeetyeet",
    "email":"yeet@yeet.dev"
    }`))

	endpoints.signup(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusCreated)
	conn.AssertExpectations(t)
}

func TestSignupIncorrectJson(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", mock.Anything).Run(func(args mock.Arguments) {
		t.Fail()
	}).Return(nil)
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()

	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{`))

	endpoints.signup(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupNoBody(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", mock.Anything).Run(func(args mock.Arguments) {
		t.Fail()
	})
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	endpoints.signup(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupNoUsername(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", mock.Anything).Run(func(args mock.Arguments) {
		t.Fail()
	})
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "password":"yeet",
    "email":"yeet@yeet.dev"
}`))

	endpoints.signup(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupNoPassword(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", mock.Anything).Run(func(args mock.Arguments) {
		t.Fail()
	})
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "email":"yeet@yeet.dev"
}`))

	endpoints.signup(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupNoEmail(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", mock.Anything).Run(func(args mock.Arguments) {
		t.Fail()
	})
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeet"
}`))

	endpoints.signup(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestSignupAdmin(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", mock.Anything).Run(func(args mock.Arguments) {
		t.Fail()
	}).Return(nil)
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeet",
    "email":"yeet@yeet.dev",
    "role": 1
    }`))

	endpoints.signup(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestDbError(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", mock.Anything).Return(errors.New("errortest"))
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeetyeet",
    "email":"yeet@yeet.dev"
}`))

	endpoints.signup(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusInternalServerError)
}

func TestLogin(t *testing.T) {
	conn := SQLConnectionMock{}
	hash, err := hash2.HashPassword("yeetyeet")
	assert.Nil(t, err)

	conn.On("GetUserByName", mock.Anything).Run(func(args mock.Arguments) {
		assert.Equal(t, args.Get(0), "jonathan")
	}).Return(db.User{
		Username: "jonathan",
		Password: hash,
		Email:    "yeet@yeet.dev",
		Role:     db.UserRoleID,
	}, nil)
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeetyeet"
    }`))

	endpoints.login(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	resp := w.Result()

	tp := jwt.TokenPair{}
	err = json.NewDecoder(resp.Body).Decode(&tp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tp.LoginToken)
	assert.NotEmpty(t, tp.RefreshToken)

	// Login token
	claims, err := jwt.VerifyJWT(tp.LoginToken, cfg)
	if assert.NoError(t, err) {
		assert.NotNil(t, claims)
		assert.Equal(t, "jonathan", claims.Username)
		assert.Equal(t, 0, claims.Role)
		assert.Equal(t, false, claims.Refresh)
	}

	// Refresh token
	claims, err = jwt.VerifyJWT(tp.RefreshToken, cfg)
	if assert.NoError(t, err) {
		assert.NotNil(t, claims)
		assert.Equal(t, "jonathan", claims.Username)
		assert.Equal(t, 0, claims.Role)
		assert.Equal(t, true, claims.Refresh)
	}
}

func TestLoginNoBody(t *testing.T) {
	conn := SQLConnectionMock{}

	conn.On("GetUserByName", mock.Anything).Run(func(args mock.Arguments) {
		t.Fail()
	}).Return(nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	endpoints.login(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestLoginIncorrectJson(t *testing.T) {
	conn := SQLConnectionMock{}

	conn.On("GetUserByName", mock.Anything).Run(func(args mock.Arguments) {
		t.Fail()
	}).Return(nil)
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{"))

	endpoints.login(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}

func TestLoginInvalidPassword(t *testing.T) {
	conn := SQLConnectionMock{}
	hash, err := hash2.HashPassword("woolnoo")
	assert.Nil(t, err)

	conn.On("GetUserByName", mock.Anything).Run(func(args mock.Arguments) {
		assert.Equal(t, args.Get(0), "jonathan")
	}).Return(db.User{
		Username: "jonathan",
		Password: hash,
		Email:    "yeet@yeet.dev",
		Role:     db.UserRoleID,
	}, nil)

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeet"
}`))

	endpoints.login(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusUnauthorized)
	conn.AssertExpectations(t)
}

func TestLoginInvalidUsername(t *testing.T) {
	conn := db.InitDB("inmemory")

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeetyeet",
    "email":"test@test.test"
}`))

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathanx",    
    "password":"yeetyeet"
}`))

	endpoints.signup(w1, r1)
	endpoints.login(w2, r2)

	assert.Equal(t, w1.Result().StatusCode, http.StatusCreated)
	assert.Equal(t, w2.Result().StatusCode, http.StatusUnauthorized)
}

func TestLoginUsernameMismatch(t *testing.T) {
	conn := db.InitDB("inmemory")

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeetyeet",
    "email":"test@test.test"
}`))

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathanx",    
    "password":"yeetyeet1",
    "email":"test@test.test"
}`))

	w3 := httptest.NewRecorder()
	r3 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathanx",    
    "password":"yeetyeet"
}`))

	endpoints.signup(w1, r1)
	endpoints.signup(w2, r2)
	endpoints.login(w3, r3)

	assert.Equal(t, w1.Result().StatusCode, http.StatusCreated)
	assert.Equal(t, w2.Result().StatusCode, http.StatusCreated)
	assert.Equal(t, w3.Result().StatusCode, http.StatusUnauthorized)
}

func TestSignupExists(t *testing.T) {
	conn := db.InitDB("inmemory")

	cfg := config.GetDefault()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yeetyeet",
    "email":"test@test.test"
    }`))

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
    "username":"jonathan",    
    "password":"yyeetyeet",
    "email":"test@test.test"
    }`))

	endpoints.signup(w1, r1)
	endpoints.signup(w2, r2)

	assert.Equal(t, w1.Result().StatusCode, http.StatusCreated)
	assert.Equal(t, w2.Result().StatusCode, http.StatusConflict)
}

func TestRefresh(t *testing.T) {
	u := db.User{
		Username: "victor",
	}

	conn := SQLConnectionMock{}
	conn.On("GetUserByName", "victor").Return(u, nil)
	conn.On("CountUsers", mock.Anything).Return(1, nil)

	cfg := config.GetDefault()

	ep := Endpoints{conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, true, cfg)
	assert.NoError(t, err)

	b, err := json.Marshal(jwt.TokenPair{
		RefreshToken: tkn,
	})

	// Create httptest recorder and request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(b))

	// Call the refresh endpoint
	ep.refresh(w, r)

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
	c, err := jwt.VerifyJWT(tp.LoginToken, cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, c)
	assert.Equal(t, u.Username, c.Username)
}
