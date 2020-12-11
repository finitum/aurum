package aurum

import (
	"strings"

	"github.com/finitum/aurum/pkg/api"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"github.com/finitum/aurum/pkg/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type RemoteClient struct {
	url string
	pk  ecc.PublicKey
}

func (a *RemoteClient) GetUrl() string {
	return a.url
}

func NewRemoteClient(url string) (*RemoteClient, error) {
	if !strings.HasPrefix(url, "https") {
		log.Warnf("[aurum] using insecure url %s, security can not be guaranteed!", url)
	}

	pkr, err := api.GetPublicKey(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting public key")
	}

	if pkr == nil {
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

func (a *RemoteClient) Login(username, password string) (*jwt.TokenPair, error) {
	tp, err := api.Login(a.url, models.User{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "login request failed")
	}

	return tp, nil
}

func (a *RemoteClient) Register(username, password, email string) error {
	return errors.Wrap(api.SignUp(a.url, models.User{
		Username: username,
		Password: password,
		Email:    email,
	}), "signup request failed")
}

func (a *RemoteClient) Verify(token string) (*jwt.Claims, error) {
	return jwt.VerifyJWT(token, a.pk)
}

func (a *RemoteClient) Refresh(tp *jwt.TokenPair) error {
	return errors.Wrap(api.Refresh(a.url, tp), "refresh client request failed")
}

func (a *RemoteClient) GetUserInfo(tp *jwt.TokenPair) (*models.User, error) {
	user, err := api.GetUser(a.url, tp)
	return user, errors.Wrap(err, "get user client request failed")
}

func (a *RemoteClient) UpdateUser(tp *jwt.TokenPair, user *models.User) (*models.User, error) {
	user, err := api.UpdateUser(a.url, tp, user)
	return user, errors.Wrap(err, "update user api request failed")
}

func (a *RemoteClient) AddGroup(tp *jwt.TokenPair, group *models.Group) error {
	err := api.AddGroup(a.url, tp, group)
	return errors.Wrap(err, "add group api request failed")
}

func (a *RemoteClient) RemoveGroup(tp *jwt.TokenPair, group string) error {
	err := api.RemoveGroup(a.url, tp, group)
	return errors.Wrap(err, "remove group api request failed")
}

func (a *RemoteClient) GetAccess(group, user string) (models.AccessStatus, error) {
	access, err := api.GetAccess(a.url, group, user)
	return access, errors.Wrap(err, "GetAccess api request failed")
}

func (a *RemoteClient) SetAccess(tp *jwt.TokenPair, access models.AccessStatus) error {
	err := api.SetAccess(a.url, tp, access)
	return errors.Wrap(err, "SetAccess api request failed")
}

func (a *RemoteClient) AddUserToGroup(tp *jwt.TokenPair, user, group string) error {
	err := api.AddUserToGroup(a.url, tp, user, group)
	return errors.Wrap(err, "AddUserToGroup api request failed")
}

func (a *RemoteClient) RemoveUserFromGroup(tp *jwt.TokenPair, user, group string) error {
	err := api.RemoveUserFromGroup(a.url, tp, user, group)
	return errors.Wrap(err, "RemoveUserFromGroup api request failed")
}

func (a *RemoteClient) GetGroupsForUser(tp *jwt.TokenPair, user string) ([]models.GroupWithRole, error) {
	groups, err := api.GetGroupsForUser(a.url, tp, user)
	if err != nil {
		return nil, errors.Wrap(err, "GetGroupsForUser api request failed")
	}
	return groups, nil
}

func (a *RemoteClient) GetUsers(tp *jwt.TokenPair) ([]models.User, error) {
	users, err := api.GetUsers(a.url, tp)
	if err != nil {
		return nil, errors.Wrap(err, "GetUsers api request failed")
	}
	return users, nil
}

func (a *RemoteClient) GetGroups() ([]models.Group, error) {
	groups, err := api.GetGroups(a.url)
	if err != nil {
		return nil, errors.Wrap(err, "GetGroups api request failed")
	}
	return groups, nil
}