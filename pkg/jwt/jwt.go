package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"github.com/finitum/aurum/pkg/models"
	"github.com/google/uuid"
	"time"
)

type Claims struct {
	Username string
	Role     models.Role
	Refresh  bool
	jwt.StandardClaims
}

type TokenPair struct {
	LoginToken   string `json:"login_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func GenerateJWT(user models.User, refresh bool, key ecc.SecretKey) (string, error) {
	// expirationTime := time.Now().Add(time.Hour)
	var expirationTime time.Time

	if refresh {
		expirationTime = time.Now().AddDate(0, 3, 0)
	} else {
		expirationTime = time.Now().Add(time.Minute * 15)
	}

	now := time.Now()
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		Refresh:  refresh,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix seconds
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Id:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(&ecc.SigningMethodEdDSA{}, claims)

	return token.SignedString(key)
}

func GenerateJWTPair(user models.User, key ecc.SecretKey) (TokenPair, error) {
	login, err := GenerateJWT(user, false, key)
	if err != nil {
		return TokenPair{}, err
	}

	refresh, err := GenerateJWT(user, true, key)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{login, refresh}, nil
}

func VerifyJWT(token string, key ecc.PublicKey) (*Claims, error) {
	claims := &Claims{}

	if _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*ecc.SigningMethodEdDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	}); err != nil {
		return nil, err
	}

	return claims, nil
}
