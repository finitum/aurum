package jwt

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"github.com/finitum/aurum/pkg/models"
	tassert "github.com/stretchr/testify/assert"
)

func TestGenerateJWTSimple(t *testing.T) {
	assert := tassert.New(t)
	cfg := config.EphemeralConfig()

	testUser := models.User{
		Username: "User",
	}

	token, err := GenerateJWT(testUser.Username, false, cfg.SecretKey)
	assert.Nil(err)
	assert.NotNil(token)

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		_, err := token.Method.(*ecc.SigningMethodEdDSA)
		assert.NotNil(err)

		return cfg.PublicKey, nil
	})

	// Assert token is valid
	assert.Nil(err)
	assert.NotNil(tkn)
	assert.True(tkn.Valid)

	// Assert token contains what we expect it to
	assert.Equal(claims.Username, testUser.Username)
	assert.WithinDuration(time.Now().Add(time.Minute*15), time.Unix(claims.ExpiresAt, 0), time.Minute)
	assert.WithinDuration(time.Now(), time.Unix(claims.IssuedAt, 0), time.Minute)
}

func TestVerifyTokenSimple(t *testing.T) {
	assert := tassert.New(t)
	cfg := config.EphemeralConfig()

	testUser := models.User{
		Username: "User",
	}

	token, err := GenerateJWT(testUser.Username, false, cfg.SecretKey)
	assert.Nil(err)
	assert.NotNil(token)

	claims, err := VerifyJWT(token, cfg.PublicKey)

	assert.Nil(err)
	assert.NotNil(claims)
	assert.Equal(claims.Username, testUser.Username)
}

func TestTokenPair(t *testing.T) {
	assert := tassert.New(t)
	cfg := config.EphemeralConfig()

	testUser := models.User{
		Username: "User",
	}
	tp, err := GenerateJWTPair(testUser.Username, cfg.SecretKey)
	assert.NotNil(tp)
	assert.Nil(err)

	claims, err := VerifyJWT(tp.LoginToken, cfg.PublicKey)
	assert.Nil(err)
	assert.NotNil(claims)
	assert.Equal(claims.Refresh, false)

	claims, err = VerifyJWT(tp.RefreshToken, cfg.PublicKey)
	assert.Nil(err)
	assert.NotNil(claims)
	assert.Equal(claims.Refresh, true)
}

func TestExpiredToken(t *testing.T) {
	cfg := config.EphemeralConfig()
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

	token := jwt.NewWithClaims(&ecc.SigningMethodEdDSA{}, claims)
	tokenString, err := token.SignedString(cfg.SecretKey)
	assert.Nil(err)
	assert.NotNil(tokenString)

	_, err = VerifyJWT(tokenString, cfg.PublicKey)
	assert.Error(err)
}

func TestWrongSigningMethod(t *testing.T) {
	assert := tassert.New(t)
	cfg := config.EphemeralConfig()

	tokenString := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0=.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ==."

	_, err := VerifyJWT(tokenString, cfg.PublicKey)
	assert.Error(err)
}

func TestInvalidJWT(t *testing.T) {
	assert := tassert.New(t)
	cfg := config.EphemeralConfig()

	tokenString := "This is clearly an invalid JWT Token"

	_, err := VerifyJWT(tokenString, cfg.PublicKey)
	assert.Error(err)
}

func TestJWT(t *testing.T) {
	assert := tassert.New(t)
	cfg := config.EphemeralConfig()

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	_, err := VerifyJWT(tokenString, cfg.PublicKey)
	assert.Error(err)
}
