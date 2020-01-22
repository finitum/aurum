package jwt

import (
	"aurum/config"
	"aurum/db"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	Username string
	Role     int
	Refresh  bool
	jwt.StandardClaims
}

type TokenPair struct {
	LoginToken   string `json:"login_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func GenerateJWT(user *db.User, refresh bool, cfg *config.Config) (string, error) {
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
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(cfg.JWTKey)
}

func GenerateJWTPair(user *db.User, cfg *config.Config) (TokenPair, error) {
	login, erra := GenerateJWT(user, false, cfg)
	refresh, errb := GenerateJWT(user, true, cfg)

	if erra != nil {
		errb = erra
	}

	return TokenPair{login, refresh}, errb
}

func VerifyJWT(token string, cfg *config.Config) (*Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return cfg.JWTKey, nil
	})

	return claims, err
}
