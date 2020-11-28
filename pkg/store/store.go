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
	// CreateGroup creates a new group in the database.
	// Group names and ids must be unique.
	CreateGroup(ctx context.Context, group models.Group) error

	// RemoveGroup removes an group from the database based
	// on it's name.
	RemoveGroup(ctx context.Context, name string) error

	// GetGroup retrieves an group based on it name.
	GetGroup(ctx context.Context, name string) (*models.Group, error)

	// GetGroups lists all groups.
	GetGroups(ctx context.Context) ([]models.Group, error)

	// GetGroupsForUser lists all groups a user has a specified role in.
	GetGroupsForUser(ctx context.Context, name string) ([]models.GroupWithRole, error)

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

	// AddUserToGroup links a user to an group with a given role.
	// This role is the role the user has within this group.
	AddGroupToUser(ctx context.Context, user string, name string, role models.Role) error

	// RemoveUserFromGroup removes the link between a user and an group.
	RemoveGroupFromUser(ctx context.Context, group string, user string) error

	// GetGroupRole retrieves the role a user has within an group
	GetGroupRole(ctx context.Context, group string, user string) (models.Role, error)

	// SetGroupRole changes the role of a user within an group.
	SetGroupRole(ctx context.Context, group string, user string, role models.Role) error

	// CountUsers counts the number of users currently in the database
	CountUsers(ctx context.Context) (int, error)
}
