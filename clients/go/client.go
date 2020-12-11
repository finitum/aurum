package aurum

import (
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
)

type Client interface {
	// User
	Login(username, password string) (*jwt.TokenPair, error)
	Register(username, password, email string) error
	Verify(token string) (*jwt.Claims, error)
	GetUserInfo(tp *jwt.TokenPair) (*models.User, error)
	UpdateUser(tp *jwt.TokenPair, user *models.User) (*models.User, error)

	GetUsers(tp *jwt.TokenPair) ([]models.User, error)

	// Group
	AddGroup(tp *jwt.TokenPair, group *models.Group) error
	RemoveGroup(tp *jwt.TokenPair, group string) error
	GetAccess(group, user string) (models.AccessStatus, error)
	SetAccess(tp *jwt.TokenPair, access models.AccessStatus) error
	AddUserToGroup(tp *jwt.TokenPair, user, group string) error
	RemoveUserFromGroup(tp *jwt.TokenPair, user, group string) error

	GetGroupsForUser(tp *jwt.TokenPair, user string) ([]models.GroupWithRole, error)
	GetGroups() ([]models.Group, error)

	GetUrl() string
}
