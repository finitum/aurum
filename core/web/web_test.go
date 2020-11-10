package web

import (
	"github.com/finitum/aurum/core/config"
	"github.com/finitum/aurum/internal/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HandlerMock struct {
	mock.Mock
}

func (h *HandlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}

func TestAuthenticateRequest(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.EphemeralConfig()

	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	conn.On("GetUserByName", "victor").Return(u, nil)

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	hm := &HandlerMock{}
	hm.On("ServeHTTP", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		// assert.NotNil(t, )
		// claims := r.Context().Value(contextKeyClaims).(jwt.Claims)
		// assert.Equal(t, u.Username, claims.Username)
		req := args.Get(1).(*http.Request)
		claims := req.Context().Value(contextKeyClaims).(*jwt.Claims)
		user := req.Context().Value(contextKeyUser).(*models.User)
		assert.Equal(t, u.Username, claims.Username)
		assert.Equal(t, u.Username, user.Username)
	})

	// Run the method
	handler := endpoints.authenticationMiddleware(hm)
	handler.ServeHTTP(w, r)

	// Assert mocks
	conn.AssertExpectations(t)
	hm.AssertExpectations(t)
}

func TestAuthenticateRequestRefreshToken(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.EphemeralConfig()
	endpoints := Endpoints{&conn, cfg}

	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  false,
	}

	tkn, err := jwt.GenerateJWT(&u, true, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	hm := &HandlerMock{}

	handler := endpoints.authenticationMiddleware(hm)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	hm.AssertExpectations(t)
	conn.AssertExpectations(t)
}

func TestAuthenticateRequestNoToken(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.EphemeralConfig()
	endpoints := Endpoints{&conn, cfg}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer ")

	hm := &HandlerMock{}

	handler := endpoints.authenticationMiddleware(hm)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

	hm.AssertExpectations(t)
	conn.AssertExpectations(t)
}

func TestAuthenticateRequestNoHeader(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.EphemeralConfig()
	endpoints := Endpoints{&conn, cfg}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Del("Authorization") // Make sure header doesn't exist

	hm := &HandlerMock{}

	handler := endpoints.authenticationMiddleware(hm)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

	hm.AssertExpectations(t)
	conn.AssertExpectations(t)
}

func TestAuthenticateRequestInvalidToken(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.EphemeralConfig()
	endpoints := Endpoints{&conn, cfg}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer "+"invalidtkn=")

	hm := &HandlerMock{}

	handler := endpoints.authenticationMiddleware(hm)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

	hm.AssertExpectations(t)
	conn.AssertExpectations(t)
}

func TestAuthenticateRequestBlockedUser(t *testing.T) {
	conn := SQLConnectionMock{}
	cfg := config.EphemeralConfig()

	u := models.User{
		Username: "victor",
		Email:    "victor@example.com",
		Role:     models.UserRoleID,
		Blocked:  true,
	}

	conn.On("GetUserByName", "victor").Return(u, nil)

	endpoints := Endpoints{&conn, cfg}

	tkn, err := jwt.GenerateJWT(&u, false, cfg.SecretKey)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Authorization", "Bearer "+tkn)

	hm := &HandlerMock{}

	handler := endpoints.authenticationMiddleware(hm)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

	hm.AssertExpectations(t)
	conn.AssertExpectations(t)
}

func TestCORSMiddlewareOptions(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodOptions, "/", nil)

	hm := &HandlerMock{}

	handler := corsMiddleware(hm)

	handler.ServeHTTP(w, r)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, OPTIONS, PUT, DELETE", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin, Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))

	hm.AssertExpectations(t)
}

func TestCORSMiddlewareNotOptions(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	hm := &HandlerMock{}
	hm.On("ServeHTTP", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		w2 := args.Get(0).(http.ResponseWriter)
		assert.Equal(t, "*", w2.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, OPTIONS, PUT, DELETE", w2.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Origin, Content-Type, Authorization", w2.Header().Get("Access-Control-Allow-Headers"))
	})

	handler := corsMiddleware(hm)

	handler.ServeHTTP(w, r)

	hm.AssertExpectations(t)
}
