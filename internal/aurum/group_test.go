package aurum

import (
	"context"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store/mock_store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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
	ms.EXPECT().CreateGroup(gomock.Any(), groupL)

	// SUT
	err = au.AddGroup(ctx, token, group)
	assert.NoError(t, err)
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
	if !registration {
		ms.EXPECT().GetGroupRole(gomock.Any(), groupL.Name, username).Return(models.RoleAdmin, nil)
		ms.EXPECT().AddGroupToUser(gomock.Any(), username, groupL.Name, models.RoleAdmin)
	} else {
		ms.EXPECT().AddGroupToUser(gomock.Any(), username, groupL.Name, models.RoleUser)
	}

	// SUT
	err = au.AddUserToGroup(ctx, token, username, groupL.Name, models.RoleAdmin)
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
