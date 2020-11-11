package store

import (
	"context"
	"errors"
	"github.com/finitum/aurum/pkg/models"
	"github.com/google/uuid"
)

var ErrExists = errors.New("already exists")
var ErrNotExists = errors.New("doesn't exist")

type AurumStore interface {
	// CreateApplication creates a new application in the database.
	// Application names and ids must be unique.
	CreateApplication(ctx context.Context, app *models.Application) error

	// RemoveApplication removes an application from the database based
	// on it's appId.
	RemoveApplication(ctx context.Context, appId uuid.UUID) error

	// GetApplication retrieves an application based on it appId.
	GetApplication(ctx context.Context, appId uuid.UUID) (*models.Application, error)

	// GetApplications lists all applications.
	GetApplications(ctx context.Context) ([]models.Application, error)

	// CreateUser creates a new user in the database.
	// User names and ids must be unique
	CreateUser(ctx context.Context, user *models.User) error

	// RemoveUser removes a user from the database
	RemoveUser(ctx context.Context, user string) error

	// GetUser retrieves a user from the database based on it's
	// user id.
	GetUser(ctx context.Context, user string) (*models.User, error)

	// GetUsers lists all users.
	GetUsers(ctx context.Context) ([]models.User, error)

	// AddUserToApplication links a user to an application with a given role.
	// This role is the role the user has within this application.
	AddApplicationToUser(ctx context.Context, user string, appId uuid.UUID, role models.Role) error

	// RemoveUserFromApplication removes the link between a user and an application.
	RemoveApplicationFromUser(ctx context.Context, user string, appId uuid.UUID) error

	// GetApplicationRole retrieves the role a user has within an application
	GetApplicationRole(ctx context.Context, user string, appId uuid.UUID) (models.Role, error)

	// SetApplicationRole changes the role of a user within an application.
	SetApplicationRole(ctx context.Context, user string, appId uuid.UUID, role models.Role) error
}

