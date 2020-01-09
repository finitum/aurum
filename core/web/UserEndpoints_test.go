package web

import (
	"aurum/config"
	"aurum/db"
	"aurum/util"
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

func TestSignup(t *testing.T) {
	conn := SQLConnectionMock{}
	conn.On("CreateUser", db.User{
		Username: "jonathan",
		Password: "yeet",
		Email:    "yeet@yeet.dev",
	}).Return(nil)

	cfg := new(config.Builder).SetDefault().Build()

	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
	"username":"jonathan",
	"password":"yeet",
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

	cfg := new(config.Builder).SetDefault().Build()

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

	cfg := new(config.Builder).SetDefault().Build()
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

	cfg := new(config.Builder).SetDefault().Build()
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

	cfg := new(config.Builder).SetDefault().Build()
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

	cfg := new(config.Builder).SetDefault().Build()
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

	cfg := new(config.Builder).SetDefault().Build()
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

	cfg := new(config.Builder).SetDefault().Build()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
	"username":"jonathan",	
	"password":"yeet",
	"email":"yeet@yeet.dev"
}`))

	endpoints.signup(w, r)

	assert.Equal(t, w.Result().StatusCode, http.StatusInternalServerError)
}

func TestLogin(t *testing.T) {
	conn := SQLConnectionMock{}
	hash, err := util.HashPassword("yeet")
	assert.Nil(t, err)

	conn.On("GetUserByName", mock.Anything).Run(func(args mock.Arguments) {
		assert.Equal(t, args.Get(0), "jonathan")
	}).Return(db.User{
		Username: "jonathan",
		Password: hash,
		Email:    "yeet@yeet.dev",
		Role:     0,
	}, nil)

	cfg := new(config.Builder).SetDefault().Build()
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

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestLoginNoBody(t *testing.T) {
	conn := SQLConnectionMock{}

	conn.On("GetUserByName", mock.Anything).Run(func(args mock.Arguments) {
		t.Fail()
	}).Return(nil)

	cfg := new(config.Builder).SetDefault().Build()
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

	cfg := new(config.Builder).SetDefault().Build()
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
	hash, err := util.HashPassword("woolnoo")
	assert.Nil(t, err)

	conn.On("GetUserByName", mock.Anything).Run(func(args mock.Arguments) {
		assert.Equal(t, args.Get(0), "jonathan")
	}).Return(db.User{
		Username: "jonathan",
		Password: hash,
		Email:    "yeet@yeet.dev",
		Role:     0,
	}, nil)

	cfg := new(config.Builder).SetDefault().Build()
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

	cfg := new(config.Builder).SetDefault().Build()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
	"username":"jonathan",	
	"password":"yeet",
	"email":"test@test.test"
}`))

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
	"username":"jonathanx",	
	"password":"yeet"
}`))

	endpoints.signup(w1, r1)
	endpoints.login(w2, r2)

	assert.Equal(t, w1.Result().StatusCode, http.StatusCreated)
	assert.Equal(t, w2.Result().StatusCode, http.StatusUnauthorized)
}

func TestLoginUsernameMismatch(t *testing.T) {
	conn := db.InitDB("inmemory")

	cfg := new(config.Builder).SetDefault().Build()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
	"username":"jonathan",	
	"password":"yeet",
	"email":"test@test.test"
}`))

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
	"username":"jonathanx",	
	"password":"yeet1",
	"email":"test@test.test"
}`))

	w3 := httptest.NewRecorder()
	r3 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
	"username":"jonathanx",	
	"password":"yeet"
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

	cfg := new(config.Builder).SetDefault().Build()
	endpoints := Endpoints{
		conn:   conn,
		config: cfg,
	}

	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
	"username":"jonathan",	
	"password":"yeet",
	"email":"test@test.test"
}`))

	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
	"username":"jonathan",	
	"password":"yeet",
	"email":"test@test.test"
}`))

	endpoints.signup(w1, r1)
	endpoints.signup(w2, r2)

	assert.Equal(t, w1.Result().StatusCode, http.StatusCreated)
	assert.Equal(t, w2.Result().StatusCode, http.StatusConflict)
}
