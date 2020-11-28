package store

import (
	"context"
	"errors"
	"github.com/finitum/aurum/pkg/models"
)

var ErrExists = errors.New("already exists")
var ErrNotExists = errors.New("doesn't exist")

//go:generate mockgen -destination mock_store/mock_store.go . AurumStore

type AurumStore interface {
	// CreateUser creates a new user in the database.
	// User names and ids must be unique
	CreateUser(ctx context.Context, user models.User) error

	// RemoveUser removes a user from the database
	RemoveUser(ctx context.Context, user string) error

	// GetUser retrieves a user from the database based on it's
	// user id.
	GetUser(ctx context.Context, user string) (models.User, error)

	// GetUsers lists all users.
	GetUsers(ctx context.Context) ([]models.User, error)

	// SetUser updates a users info in the database.
	// User names and ids must be the same
	SetUser(ctx context.Context, user models.User) (models.User, error)

	// CountUsers counts the number of users currently in the database
	CountUsers(ctx context.Context) (int, error)
}
