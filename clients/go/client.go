package aurum

import (
	"context"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
)

type Client interface {
	// User
	Login(ctx context.Context, username, password string) (*jwt.TokenPair, error)
	Register(ctx context.Context, username, password, email string) error
	Verify(ctx context.Context, token string) (*jwt.Claims, error)
	GetUserInfo(ctx context.Context, tp *jwt.TokenPair) (*models.User, error)
	UpdateUser(ctx context.Context,tp *jwt.TokenPair, user *models.User) (*models.User, error)

	// Application
	AddApplication(ctx context.Context,tp *jwt.TokenPair, app *models.Application) error
	RemoveApplication(ctx context.Context,tp *jwt.TokenPair, app string) error
	GetAccess(ctx context.Context,app, user string) (models.AccessStatus, error)
	SetAccess(ctx context.Context,tp *jwt.TokenPair, access models.AccessStatus) error
	AddUserToApplication(ctx context.Context,tp *jwt.TokenPair, user, app string) error
	RemoveUserFromApplication(ctx context.Context,tp *jwt.TokenPair, user, app string) error
}
