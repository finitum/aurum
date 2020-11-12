package aurum

import (
	"context"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/pkg/errors"
	"strings"
)

func (au Aurum) AddApplication(ctx context.Context, token string, application models.Application) error {
	role, _, err := au.checkTokenAndRole(ctx, token, AurumName)
	if err != nil {
		return err
	}

	if role < models.RoleAdmin {
		return ErrUnauthorized
	}

	application.Name = strings.ToLower(application.Name)
	return au.db.CreateApplication(ctx, application)
}

func (au Aurum) RemoveApplication(ctx context.Context, token, app string) error {
	app = strings.ToLower(app)

	role, _, err := au.checkTokenAndRole(ctx, token, app)
	if err != nil {
		return err
	}

	if role < models.RoleAdmin {
		return ErrUnauthorized
	}

	if app == strings.ToLower(AurumName) {
		return errors.Errorf("Can't remove application named %s", AurumName)
	}

	return au.db.RemoveApplication(ctx, app)
}

func (au Aurum) AddUserToApplication(ctx context.Context, token, username, appName string, role models.Role) error {
	appName = strings.ToLower(appName)

	app, err := au.db.GetApplication(ctx, appName)
	if err != nil {
		return errors.Wrap(err, "getting application")
	}

	if !app.AllowRegistration {
		role, _, err := au.checkTokenAndRole(ctx, token, appName)
		if err != nil {
			return err
		}

		if role < models.RoleAdmin {
			return ErrUnauthorized
		}

	} else {
		claims, err := au.checkToken(token)
		if err != nil {
			return err
		}

		username = claims.Username
	}

	return au.db.AddApplicationToUser(ctx, username, appName, role)
}

// GetAccess determines if a user is allowed access to a certain application
func (au Aurum) GetAccess(ctx context.Context, user, name string) (models.AccessStatus, error) {
	role, err := au.db.GetApplicationRole(ctx, name, user)

	if err == store.ErrNotExists {
		return models.AccessStatus{
			ApplicationName: name,
			Username:        user,
			AllowedAccess:   false,
		}, nil

	} else if err != nil {
		return models.AccessStatus{}, err
	}

	return models.AccessStatus{
		ApplicationName: name,
		Username:        user,
		AllowedAccess:   true,
		Role:            role,
	}, nil
}

func (au Aurum) SetAccess(ctx context.Context, token, app string, role models.Role) error {
	app = strings.ToLower(app)

	role, claims, err := au.checkTokenAndRole(ctx, token, app)
	if err != nil {
		return err
	}

	if role < models.RoleAdmin {
		return ErrUnauthorized
	}

	return au.db.SetApplicationRole(ctx, claims.Username, app, role)
}

func (au Aurum) RemoveUserFromApplication(ctx context.Context, token, app string) error {
	app = strings.ToLower(app)

	role, claims, err := au.checkTokenAndRole(ctx, token, app)
	if err != nil {
		return err
	}

	if role < models.RoleAdmin {
		return ErrUnauthorized
	}

	return au.db.RemoveApplicationFromUser(ctx, claims.Username, app)
}
