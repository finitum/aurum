package aurum

import (
	"context"
	"strings"

	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/pkg/errors"
)

func (au Aurum) AddGroup(ctx context.Context, token string, group models.Group) error {
	role, claims, err := au.checkTokenAndRole(ctx, token, AurumName)
	if err != nil {
		return err
	}

	if role < models.RoleAdmin {
		return ErrUnauthorized
	}

	group.Name = strings.ToLower(group.Name)
	if err := au.db.CreateGroup(ctx, group); err != nil {
		return err
	}

	return au.db.AddGroupToUser(ctx, claims.Username, group.Name, models.RoleAdmin)
}

func (au Aurum) RemoveGroup(ctx context.Context, token, group string) error {
	group = strings.ToLower(group)

	role, _, err := au.checkTokenAndRole(ctx, token, group)
	if err != nil {
		return err
	}

	if role < models.RoleAdmin {
		return ErrUnauthorized
	}

	if group == strings.ToLower(AurumName) {
		return errors.Errorf("Can't remove group named %s", AurumName)
	}

	return au.db.RemoveGroup(ctx, group)
}

// GetAccess determines if a user is allowed access to a certain group
func (au Aurum) GetAccess(ctx context.Context, user, group string) (models.AccessStatus, error) {
	group = strings.ToLower(group)
	role, err := au.db.GetGroupRole(ctx, group, user)

	if err == store.ErrNotExists {
		return models.AccessStatus{
			GroupName:     group,
			Username:      user,
			AllowedAccess: false,
		}, nil
	} else if err != nil {
		return models.AccessStatus{}, err
	}

	return models.AccessStatus{
		GroupName:     group,
		Username:      user,
		AllowedAccess: true,
		Role:          role,
	}, nil
}

func (au Aurum) SetAccess(ctx context.Context, token, group, username string, targetRole models.Role) error {
	group = strings.ToLower(group)

	role, _, err := au.checkTokenAndRole(ctx, token, group)
	if err != nil {
		return err
	}

	if role < models.RoleAdmin {
		return ErrUnauthorized
	}

	return au.db.SetGroupRole(ctx, group, username, targetRole)
}

func (au Aurum) AddUserToGroup(ctx context.Context, token, username, groupName string, wanted models.Role) error {
	groupName = strings.ToLower(groupName)

	group, err := au.db.GetGroup(ctx, groupName)
	if err != nil {
		return errors.Wrap(err, "getting group")
	}

	role, claims, err := au.checkTokenAndRole(ctx, token, group.Name)
	if err != nil {
		return errors.Wrap(err, "getting token and role")
	}

	if role == models.RoleAdmin {
		return au.db.AddGroupToUser(ctx, username, group.Name, wanted)
	} else if wanted > models.RoleUser || username != claims.Username || !group.AllowRegistration {
		return ErrUnauthorized
	}

	return au.db.AddGroupToUser(ctx, claims.Username, group.Name, models.RoleUser)
}

func (au Aurum) RemoveUserFromGroup(ctx context.Context, token, target, group string) error {
	group = strings.ToLower(group)

	role, claims, err := au.checkTokenAndRole(ctx, token, group)
	if err != nil {
		return err
	}

	if role < models.RoleAdmin && target != claims.Username {
		return ErrUnauthorized
	}

	return au.db.RemoveGroupFromUser(ctx, group, target)
}

func (au Aurum) GetGroupsForUser(ctx context.Context, token, user string) ([]models.GroupWithRole, error) {
	claims, err := au.checkToken(token)
	if err != nil {
		return nil, err
	}

	if user != claims.Username {
		role, err := au.checkRole(ctx, claims, AurumName)
		if err != nil {
			return nil, err
		}

		if role < models.RoleAdmin {
			// Only admins may see groups for other users
			return nil, ErrUnauthorized
		}
	}

	return au.db.GetGroupsForUser(ctx, user)
}


func (au Aurum) GetGroups(ctx context.Context) ([]models.Group, error) {
	groups, err := au.db.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
