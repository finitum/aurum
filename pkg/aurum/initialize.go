package aurum

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/finitum/aurum/internal/hash"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const adminUsername = "admin"
const Aurum = "Aurum"

func Initialize(ctx context.Context, db store.AurumStore) error {
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

	if err := db.CreateApplication(ctx, models.Application{
		Name: Aurum,
	}); err != nil {
		return errors.Wrap(err, "create initial application")
	}

	if err := db.AddApplicationToUser(ctx, adminUsername, Aurum, models.RoleAdmin); err != nil {
		return errors.Wrap(err, "add initial user to Aurum application")
	}

	return nil
}
