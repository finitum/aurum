package aurum

import (
	"context"
	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/pkg/errors"
)

type EmbeddedClient struct {
	au aurum.Aurum
}

func NewEmbeddedClient(ctx context.Context, db store.AurumStore, cfg *config.Config) (EmbeddedClient, error) {
	au, err := aurum.New(ctx, db, cfg)

	if err != nil {
		return EmbeddedClient{}, errors.Wrap(err, "failed creating aurum client")
	}

	return EmbeddedClient{au}, nil
}

func (e *EmbeddedClient) Login(ctx context.Context, username, password string) (*jwt.TokenPair, error) {
	tp, err := e.au.Login(ctx, models.User{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &tp, nil
}

func (e *EmbeddedClient) Register(ctx context.Context, username, password, email string) error {
	return e.au.SignUp(ctx, models.User{
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
