package goclient

import (
	"github.com/finitum/aurum/clients/goclient/requests"
	"github.com/finitum/aurum/internal/jwt/ecc"
	"github.com/pkg/errors"
)

type Aurum struct {
	url string
	pk ecc.PublicKey
}

func Connect(url string) (*Aurum,error) {
	pkr, err := requests.GetPublicKey(url)
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

	return  &Aurum{url, pk}, nil
}

func (a *Aurum) Login(username, password string) (string, error) {
	panic("not implemented!")
}

func (a *Aurum) Register(username, password string) (string, error) {
	panic("not implemented!")
}

func (a *Aurum) Verify(token string) (bool, error) {
	panic("not implemented!")
}
