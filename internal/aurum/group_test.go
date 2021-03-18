package aurum

import (
	"context"
	"github.com/finitum/aurum/pkg/store"
	"github.com/pkg/errors"
	"strings"
	"testing"

	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store/mock_store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAurum_AddGroup(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	group := models.Group{
		Name:              "NotAurum",
		AllowRegistration: true,
	}

	groupL := group
	groupL.Name = strings.ToLower(group.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), strings.ToLower(AurumName), "bob").Return(models.RoleAdmin, nil)
	ms.EXPECT().GetGroup(gomock.Any(), groupL.Name).Return(nil, store.ErrNotExists)
	ms.EXPECT().CreateGroup(gomock.Any(), groupL)
	ms.EXPECT().AddGroupToUser(gomock.Any(), "bob", groupL.Name, models.RoleAdmin)

	// SUT
	err = au.AddGroup(ctx, token, group)
	assert.NoError(t, err)
}

func TestAurum_AddGroup_Duplicate(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	group := models.Group{
		Name:              "Aurum",
		AllowRegistration: true,
	}

	groupL := group
	groupL.Name = strings.ToLower(group.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), strings.ToLower(AurumName), "bob").Return(models.RoleAdmin, nil)
	ms.EXPECT().GetGroup(gomock.Any(), groupL.Name).Return(&group, nil)

	// SUT
	err = au.AddGroup(ctx, token, group)
	assert.Error(t, err)
}

func TestAurum_AddGroup_Unauthorized(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	group := models.Group{
		Name:              "NotAurum",
		AllowRegistration: true,
	}

	groupL := group
	groupL.Name = strings.ToLower(group.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), strings.ToLower(AurumName), "bob").Return(models.RoleUser, nil)

	// SUT
	err = au.AddGroup(ctx, token, group)
	assert.Error(t, err)
}

func TestAurum_RemoveGroup(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	group := models.Group{
		Name:              "NotAurum",
		AllowRegistration: true,
	}

	groupL := group
	groupL.Name = strings.ToLower(group.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL.Name, "bob").Return(models.RoleAdmin, nil)
	ms.EXPECT().RemoveGroup(gomock.Any(), groupL.Name)

	// SUT
	err = au.RemoveGroup(ctx, token, groupL.Name)
	assert.NoError(t, err)
}

func TestAurum_RemoveGroup_Unauthorized(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	group := models.Group{
		Name:              "NotAurum",
		AllowRegistration: true,
	}

	groupL := group
	groupL.Name = strings.ToLower(group.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL.Name, "bob").Return(models.RoleUser, nil)

	// SUT
	err = au.RemoveGroup(ctx, token, groupL.Name)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestAurum_RemoveGroup_NonExistent(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	group := models.Group{
		Name:              "NotAurum",
		AllowRegistration: true,
	}

	groupL := group
	groupL.Name = strings.ToLower(group.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL.Name, "bob").Return(models.Role(0), store.ErrNotExists)

	// SUT
	err = au.RemoveGroup(ctx, token, groupL.Name)
	assert.Equal(t, store.ErrNotExists, err)
}

func TestAurum_RemoveGroup_Aurum(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	group := models.Group{
		Name:              AurumName,
		AllowRegistration: true,
	}

	groupL := group
	groupL.Name = strings.ToLower(group.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL.Name, "bob").Return(models.RoleAdmin, nil)

	// SUT
	err = au.RemoveGroup(ctx, token, groupL.Name)
	assert.Error(t, err)
}

func testAddToGroupHelper(t *testing.T, registration bool) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	group := models.Group{
		Name:              "NotAurum",
		AllowRegistration: registration,
	}

	const username = "bob"

	groupL := group
	groupL.Name = strings.ToLower(group.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroup(gomock.Any(), groupL.Name).Return(&groupL, nil)
	if registration {
		ms.EXPECT().GetGroupRole(gomock.Any(), groupL.Name, username).Return(models.Role(0), nil)
	} else {
		ms.EXPECT().GetGroupRole(gomock.Any(), groupL.Name, username).Return(models.RoleAdmin, nil)
	}
	ms.EXPECT().AddGroupToUser(gomock.Any(), username, groupL.Name, models.RoleUser)

	// SUT
	err = au.AddUserToGroup(ctx, token, username, groupL.Name, models.RoleUser)
	assert.NoError(t, err)
}

func TestAurum_AddUserToGroup(t *testing.T) {
	testAddToGroupHelper(t, true)
	testAddToGroupHelper(t, false)
}

func TestAurum_GetAccess(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"

	const group = "Angroup"
	groupL := strings.ToLower(group)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.RoleAdmin, nil)

	// SUT
	au := Aurum{db: ms}
	resp, err := au.GetAccess(ctx, username, group)
	assert.NoError(t, err)

	assert.Equal(t, models.AccessStatus{
		GroupName:     groupL,
		Username:      username,
		AllowedAccess: true,
		Role:          models.RoleAdmin,
	}, resp)
}

func TestArum_GetAccess_NonExistent(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"
	const group = "AnGroup"
	groupL := strings.ToLower(group)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.Role(0), store.ErrNotExists)

	// SUT
	au := Aurum{db: ms}
	resp, err := au.GetAccess(ctx, username, group)
	assert.NoError(t, err)
	assert.Equal(t, models.AccessStatus{
		GroupName:     groupL,
		Username:      username,
		AllowedAccess: false,
		Role:          0,
	}, resp)
}

func TestArum_GetAccess_UnexpectedError(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"
	const group = "AnGroup"
	groupL := strings.ToLower(group)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.Role(0), errors.New("db rip"))

	// SUT
	au := Aurum{db: ms}
	resp, err := au.GetAccess(ctx, username, group)
	assert.Empty(t, resp)
	assert.Error(t, err)
}

func TestAurum_SetAccess(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"
	const target = "wooloo"

	const group = "Angroup"
	groupL := strings.ToLower(group)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.RoleAdmin, nil)
	ms.EXPECT().SetGroupRole(gomock.Any(), groupL, target, models.RoleUser)

	// SUT
	err = au.SetAccess(ctx, token, group, target, models.RoleUser)
	assert.NoError(t, err)
}

func TestAurum_SetAccess_Unauthorized(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"
	const target = "wooloo"

	const group = "Angroup"
	groupL := strings.ToLower(group)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.RoleUser, nil)

	// SUT
	err = au.SetAccess(ctx, token, group, target, models.RoleUser)
	assert.Equal(t, err, ErrUnauthorized)
}

func TestAurum_SetAccess_NonExistent(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"
	const target = "wooloo"

	const group = "Angroup"
	groupL := strings.ToLower(group)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.Role(0), store.ErrNotExists)

	// SUT
	err = au.SetAccess(ctx, token, group, target, models.RoleUser)
	assert.Equal(t, store.ErrNotExists, err)
}

func TestAurum_RemoveUserFromGroup(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"
	const target = "wooloo"

	const group = "Angroup"
	groupL := strings.ToLower(group)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.RoleAdmin, nil)
	ms.EXPECT().RemoveGroupFromUser(gomock.Any(), groupL, target)

	// SUT
	err = au.RemoveUserFromGroup(ctx, token, target, group)
	assert.NoError(t, err)
}

func TestAurum_RemoveUserFromGroup_Unauthorized(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"
	const target = "wooloo"

	const group = "Angroup"
	groupL := strings.ToLower(group)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.RoleUser, nil)

	// SUT
	err = au.RemoveUserFromGroup(ctx, token, target, group)
	assert.Equal(t, ErrUnauthorized, err)
}


func TestAurum_RemoveUserFromGroup_Self(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"

	const group = "Angroup"
	groupL := strings.ToLower(group)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.RoleUser, nil)
	ms.EXPECT().RemoveGroupFromUser(gomock.Any(), groupL, username)

	// SUT
	err = au.RemoveUserFromGroup(ctx, token, username, group)
	assert.NoError(t, err)
}

func TestAurum_RemoveUserFromGroup_NonExistent(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"

	const group = "Angroup"
	groupL := strings.ToLower(group)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), groupL, username).Return(models.Role(0), store.ErrNotExists)

	// SUT
	err = au.RemoveUserFromGroup(ctx, token, username, group)
	assert.Error(t, err)
}

func TestAurum_GetGroupsForUser(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()
	cfg := config.EphemeralConfig()
	ms := mock_store.NewMockAurumStore(ctrl)
	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	const username = "yeet"

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	groups := []models.GroupWithRole{
		{
			Group: models.Group{
				Name:              "aurum",
			},
			Role:  models.RoleUser,
		},
	}

	// Expect
	ms.EXPECT().GetGroupsForUser(gomock.Any(), username).Return(groups, nil)

	// SUT
	rgroups, err := au.GetGroupsForUser(ctx, token, username)
	assert.NoError(t, err)
	assert.Equal(t, groups, rgroups)
}

func TestAurum_GetGroupsForUser_Admin(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()
	cfg := config.EphemeralConfig()
	ms := mock_store.NewMockAurumStore(ctrl)
	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	const username = "yeet"
	const target = "yoinks"

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	groups := []models.GroupWithRole{
		{
			Group: models.Group{
				Name:              "aurum",
			},
			Role:  models.RoleUser,
		},
	}

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), strings.ToLower(AurumName), username).Return(models.RoleAdmin, nil)
	ms.EXPECT().GetGroupsForUser(gomock.Any(), target).Return(groups, nil)

	// SUT
	rgroups, err := au.GetGroupsForUser(ctx, token, target)
	assert.NoError(t, err)
	assert.Equal(t, groups, rgroups)
}

func TestAurum_GetGroupsForUser_Unauthorized(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()
	cfg := config.EphemeralConfig()
	ms := mock_store.NewMockAurumStore(ctrl)
	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	const username = "yeet"
	const target = "yoinks"

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupRole(gomock.Any(), strings.ToLower(AurumName), username).Return(models.RoleUser, nil)

	// SUT
	rgroups, err := au.GetGroupsForUser(ctx, token, target)
	assert.Nil(t, rgroups)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestAurum_GetGroupsForUser_NonExistent(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()
	cfg := config.EphemeralConfig()
	ms := mock_store.NewMockAurumStore(ctrl)
	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	const username = "yeet"

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetGroupsForUser(gomock.Any(), username).Return(nil, store.ErrNotExists)

	// SUT
	rgroups, err := au.GetGroupsForUser(ctx, token, username)
	assert.Nil(t, rgroups)
	assert.Equal(t, store.ErrNotExists, err)
}
