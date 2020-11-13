package aurum

import (
	"context"
	"github.com/finitum/aurum/pkg/api"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"github.com/finitum/aurum/pkg/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

type RemoteClient struct {
	url string
	pk  ecc.PublicKey
}

func NewRemoteClient(url string) (*RemoteClient, error) {
	if !strings.HasPrefix(url, "https") {
		log.Warnf("[aurum] using insecure url %s, security can not be guaranteed!", url)
	}

	pkr, err := api.GetPublicKey(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting public key")
	}

	key, err := ecc.FromPem([]byte(pkr.PublicKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing public key")
	}

	pk, ok := key.(ecc.PublicKey)
	if !ok {
		return nil, errors.New("unexpected key type")
	}

	return &RemoteClient{url, pk}, nil
}

func (a *RemoteClient) Login(_ context.Context, username, password string) (*jwt.TokenPair, error) {
	tp, err := api.Login(a.url, models.User{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "login request failed")
	}

	return tp, nil
}

func (a *RemoteClient) Register(_ context.Context, username, password, email string) error {
	return errors.Wrap(api.SignUp(a.url, models.User{
		Username: username,
		Password: password,
		Email:    email,
	}), "signup request failed")
}

func (a *RemoteClient) Verify(_ context.Context, token string) (*jwt.Claims, error) {
	return jwt.VerifyJWT(token, a.pk)
}

func (a *RemoteClient) refresh(tp *jwt.TokenPair) error {
	return errors.Wrap(api.Refresh(a.url, tp), "refresh api request failed")
}

func (a *RemoteClient) GetUserInfo(_ context.Context, tp *jwt.TokenPair) (*models.User, error) {
	user, err := api.GetUser(a.url, tp)
	return user, errors.Wrap(err, "get user api request failed")
}

func (a *RemoteClient) UpdateUser(_ context.Context, tp *jwt.TokenPair, user *models.User) (*models.User, error) {
	user, err := api.UpdateUser(a.url, tp, user)
	return user, errors.Wrap(err, "update user api request failed")
}

func (a *RemoteClient) AddApplication(_ context.Context, tp *jwt.TokenPair, app *models.Application) error {
	err := api.AddApplication(a.url, tp, app)
	return errors.Wrap(err, "add application api request failed")
}

func (a *RemoteClient) RemoveApplication(_ context.Context, tp *jwt.TokenPair, app string) error {
	err := api.RemoveApplication(a.url, tp, app)
	return errors.Wrap(err, "remove application api request failed")
}

func (a *RemoteClient) GetAccess(_ context.Context, app, user string) (models.AccessStatus, error) {
	access, err := api.GetAccess(a.url, app, user)
	return access, errors.Wrap(err, "GetAccess api request failed")
}

func (a *RemoteClient) SetAccess(_ context.Context, tp *jwt.TokenPair, access models.AccessStatus) error {
	err := api.SetAccess(a.url, tp, access)
	return errors.Wrap(err, "SetAccess api request failed")
}

func (a *RemoteClient) AddUserToApplication(_ context.Context, tp *jwt.TokenPair, user, app string) error {
	err := api.AddUserToApplication(a.url, tp, user, app)
	return errors.Wrap(err, "AddUserToApplication api request failed")
}

func (a *RemoteClient) RemoveUserFromApplication(_ context.Context, tp *jwt.TokenPair, user, app string) error {
	err := api.RemoveUserFromApplication(a.url, tp, user, app)
	return errors.Wrap(err, "RemoveUserFromApplication api request failed")
}
