package aurum

import (
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"github.com/finitum/aurum/pkg/oldapi"
	"github.com/finitum/aurum/pkg/oldmodels"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Aurum struct {
	url string
	pk  ecc.PublicKey
}

func Connect(url string) (*Aurum, error) {
	if !strings.HasPrefix(url, "https") {
		log.Warnf("[aurum] using insecure url %s, security can not be guaranteed!", url)
	}

	pkr, err := oldapi.GetPublicKey(url)
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

	return &Aurum{url, pk}, nil
}

func (a *Aurum) Login(username, password string) (*jwt.TokenPair, error) {
	tp, err := oldapi.Login(a.url, oldmodels.User{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "login request failed")
	}

	return tp, nil
}

func (a *Aurum) Register(username, password, email string) error {
	return errors.Wrap(oldapi.SignUp(a.url, oldmodels.User{
		Username: username,
		Password: password,
		Email:    email,
	}), "signup request failed")
}

func (a *Aurum) Verify(token string) (*jwt.Claims, error) {
	return jwt.VerifyJWT(token, a.pk)
}

func (a *Aurum) Refresh(tp *jwt.TokenPair) error {
	return errors.Wrap(oldapi.Refresh(a.url, tp), "refresh oldapi request failed")
}

func (a *Aurum) GetUserInfo(tp *jwt.TokenPair) (*oldmodels.User, error) {
	user, err := oldapi.GetUser(a.url, tp)
	return user, errors.Wrap(err, "get user oldapi request failed")
}
