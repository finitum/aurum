// Package aurum is the main business logic of Aurum and it contains the methods which connect the endpoints and
// the database
package aurum

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/finitum/aurum/internal/hash"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrInvalidInput = errors.New("password is too weak")
var ErrWeakPassword = errors.New("password is too weak")
var ErrUnauthorized = errors.New("unauthorized")

const adminUsername = "admin"

type Aurum struct {
	db store.AurumStore
	pk ecc.PublicKey
	sk ecc.SecretKey
}

func New(ctx context.Context, db store.AurumStore, cfg *config.Config) (Aurum, error) {
	if err := setup(ctx, db); err != nil {
		return Aurum{}, err
	}
	return Aurum{db, cfg.PublicKey, cfg.SecretKey}, nil
}

func setup(ctx context.Context, db store.AurumStore) error {
	nu, err := db.CountUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "count users")
	}

	if nu != 0 {
		return nil
	}

	log.Info("Detected first run - Initializing Aurum")

	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	if err != nil {
		return errors.Wrap(err, "random")
	}
	pass := base64.StdEncoding.EncodeToString(buf)

	log.Infof("Created initial user: \"%s\" with password \"%s\"", adminUsername, pass)

	hashed, err := hash.HashPassword(pass)
	if err != nil {
		return errors.Wrap(err, "hashing failed")
	}

	if err := db.CreateUser(ctx, models.User{
		Username: adminUsername,
		Password: hashed,
	}); err != nil {
		return errors.Wrap(err, "create initial user")
	}

	return nil
}

func (au Aurum) checkToken(token string) (*jwt.Claims, error) {
	claims, err := jwt.VerifyJWT(token, au.pk)
	if err != nil {
		return nil, ErrUnauthorized
	}

	if err := claims.Valid(); err != nil {
		return nil, ErrUnauthorized
	}

	// Refresh tokens are not allowed to be used as authentication
	if claims.Refresh {
		return nil, ErrInvalidInput
	}

	return claims, nil
}
