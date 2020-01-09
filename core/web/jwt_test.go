package web

import (
	"aurum/config"
	"aurum/db"
	"github.com/dgrijalva/jwt-go"
	tassert "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateJWTSimple(t *testing.T) {
	assert := tassert.New(t)
	cfg := new(config.Builder).SetDefault().Build()

	testUser := db.User{
		Username: "User",
		Role:     0,
	}

	token, err := GenerateJWT(&testUser, false, cfg)
	assert.Nil(err)
	assert.NotNil(token)

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		_, err := token.Method.(*jwt.SigningMethodHMAC)
		assert.NotNil(err)

		return cfg.JWTKey, nil
	})

	// Assert token is valid
	assert.Nil(err)
	assert.NotNil(tkn)
	assert.True(tkn.Valid)

	// Assert token contains what we expect it to
	assert.Equal(claims.Username, testUser.Username)
	assert.Equal(claims.Role, testUser.Role)
	assert.WithinDuration(time.Now().Add(time.Hour), time.Unix(claims.ExpiresAt, 0), time.Minute)
	assert.WithinDuration(time.Now(), time.Unix(claims.IssuedAt, 0), time.Minute)
}

func TestVerifyTokenSimple(t *testing.T) {
	assert := tassert.New(t)
	cfg := new(config.Builder).SetDefault().Build()

	testUser := db.User{
		Username: "User",
		Role:     0,
	}

	token, err := GenerateJWT(&testUser, false, cfg)
	assert.Nil(err)
	assert.NotNil(token)

	claims, err := VerifyJWT(token, cfg)

	assert.Nil(err)
	assert.NotNil(claims)
	assert.Equal(claims.Username, testUser.Username)
	assert.Equal(claims.Role, testUser.Role)
}

func TestTokenPair(t *testing.T) {
	assert := tassert.New(t)
	cfg := new(config.Builder).SetDefault().Build()

	testUser := db.User{
		Username: "User",
		Role:     0,
	}
	tp, err := GenerateJWTPair(&testUser, cfg)
	assert.NotNil(tp)
	assert.Nil(err)

	claims, err := VerifyJWT(tp.LoginToken, cfg)
	assert.Nil(err)
	assert.NotNil(claims)
	assert.Equal(claims.Refresh, false)

	claims, err = VerifyJWT(tp.RefreshToken, cfg)
	assert.Nil(err)
	assert.NotNil(claims)
	assert.Equal(claims.Refresh, true)
}

func TestExpiredToken(t *testing.T) {
	cfg := new(config.Builder).SetDefault().Build()
	assert := tassert.New(t)

	expirationTime := time.Now().Add(-(time.Hour + time.Minute))
	now := time.Now()
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(cfg.JWTKey)
	assert.Nil(err)
	assert.NotNil(tokenString)

	_, err = VerifyJWT(tokenString, cfg)
	assert.Error(err)
}

func TestWrongSigningMethod(t *testing.T) {
	assert := tassert.New(t)
	cfg := new(config.Builder).SetDefault().Build()

	tokenString := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0=.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ==."

	_, err := VerifyJWT(tokenString, cfg)
	assert.Error(err)
}

func TestInvalidJWT(t *testing.T) {
	assert := tassert.New(t)
	cfg := new(config.Builder).SetDefault().Build()

	tokenString := "This is clearly an invalid JWT Token"

	_, err := VerifyJWT(tokenString, cfg)
	assert.Error(err)
}

func TestJWT(t *testing.T) {
	assert := tassert.New(t)
	cfg := new(config.Builder).SetDefault().Build()

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	_, err := VerifyJWT(tokenString, cfg)
	assert.Error(err)
}