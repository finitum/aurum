package store

import (
	"github.com/finitum/aurum/pkg/models"
	"github.com/google/uuid"
)

type ApplicationStore interface {
	// CreateApplication creates a new application in the database.
	// Application names and ids must be unique.
	CreateApplication(*models.Application) error

	// RemoveApplication removes an application from the database based
	// on it's appId.
	RemoveApplication(appId uuid.UUID) error

	// GetApplication retrieves an application based on it appId.
	GetApplication(appId uuid.UUID) (*models.Application, error)

	// GetApplications lists all applications.
	GetApplications() ([]models.Application, error)
}

type UserStore interface {
	// CreateUser creates a new user in the database.
	// User names and ids must be unique
	CreateUser(user *models.User) error

	// RemoveUser removes a user from the database
	RemoveUser(userId uuid.UUID) error

	// GetUser retrieves a user from the database based on it's
	// user id.
	GetUser(userId uuid.UUID) (*models.User, error)

	// GetUsers lists all users.
	GetUsers() ([]models.User, error)

	// AddUserToApplication links a user to an application with a given role.
	// This role is the role the user has within this application.
	AddUserToApplication(userId uuid.UUID, appId uuid.UUID, role models.Role) error

	// RemoveUserFromApplication removes the link between a user and an application.
	RemoveUserFromApplication(userId uuid.UUID, appId uuid.UUID) error

	// SetApplicationRole sets the role changes the role of a user within an application.
	SetApplicationRole(userId uuid.UUID, appId uuid.UUID, role models.Role) error

	// GetApplicationRole retrieves the role a user has within an application
	GetApplicationRole(userId uuid.UUID, appId uuid.UUID) (models.Role, error)
}

type AurumStore interface {
	UserStore
	ApplicationStore
}
