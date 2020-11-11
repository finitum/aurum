package aurum

import (
	"context"
	"github.com/finitum/aurum/pkg/aurum"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/finitum/aurum/pkg/store/dgraph"
	"github.com/pkg/errors"
)

type EmbeddedClient struct {
	dg  store.AurumStore
	cfg config.Config
}

func NewEmbeddedClient(ctx context.Context, cfg config.Config) (*EmbeddedClient, error) {
	dg, err := dgraph.New(ctx, cfg.DgraphUrl)
	if err != nil {
		return nil, errors.Wrap(err, "dgraph connect")
	}

	return &EmbeddedClient{dg: dg, cfg: cfg}, nil
}

func (e *EmbeddedClient) Login(ctx context.Context, username, password string) (*jwt.TokenPair, error) {
	tp, err := aurum.Login(ctx, e.dg, models.User{
		Username: username,
		Password: password,
	}, e.cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	return &tp, nil
}

func (e *EmbeddedClient) Register(ctx context.Context, username, password, email string) error {
	return aurum.SignUp(ctx, e.dg, models.User{
		Username: username,
		Password: password,
		Email:    email,
	})
}

func (e *EmbeddedClient) Verify(ctx context.Context, token string) (*jwt.Claims, error) {
	panic("implement me")
}

func (e *EmbeddedClient) GetUserInfo(ctx context.Context, tp *jwt.TokenPair) (*models.User, error) {
	panic("implement me")
}
