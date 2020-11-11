package aurum

import (
	"context"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/pkg/errors"
)

func AddApplication(ctx context.Context, db store.AurumStore, application models.Application) error {
	return db.CreateApplication(ctx, application)
}

func RemoveApplication(ctx context.Context, db store.AurumStore, name string) error {
	if name == Aurum {
		return errors.Errorf("Can't remove application named %s", Aurum)
	}

	return db.RemoveApplication(ctx, name)
}

func AddUserToApplication(ctx context.Context, db store.AurumStore, username, name string, role models.Role) error {
	return db.AddApplicationToUser(ctx, username, name, role)
}

func SetRole(ctx context.Context, db store.AurumStore, username, name string, role models.Role) error {
	return db.SetApplicationRole(ctx, username, name, role)
}

func RemoveUserFromApplication(ctx context.Context, db store.AurumStore, username, name string) error {
	return db.RemoveApplicationFromUser(ctx, username, name)
}
