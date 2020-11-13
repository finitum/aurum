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

func TestAurum_AddApplication(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	app := models.Application{
		Name:              "NotAurum",
		AllowRegistration: true,
	}

	appL := app
	appL.Name = strings.ToLower(app.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetApplicationRole(gomock.Any(), strings.ToLower(AurumName), "bob").Return(models.RoleAdmin, nil)
	ms.EXPECT().CreateApplication(gomock.Any(), appL)

	// SUT
	err = au.AddApplication(ctx, token, app)
	assert.NoError(t, err)
}

func TestAurum_RemoveApplication(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	app := models.Application{
		Name:              "NotAurum",
		AllowRegistration: true,
	}

	appL := app
	appL.Name = strings.ToLower(app.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetApplicationRole(gomock.Any(), appL.Name, "bob").Return(models.RoleAdmin, nil)
	ms.EXPECT().RemoveApplication(gomock.Any(), appL.Name)

	// SUT
	err = au.RemoveApplication(ctx, token, appL.Name)
	assert.NoError(t, err)
}

func testAddToApplicationHelper(t *testing.T, registration bool) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	app := models.Application{
		Name:              "NotAurum",
		AllowRegistration: registration,
	}

	const username = "bob"

	appL := app
	appL.Name = strings.ToLower(app.Name)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetApplication(gomock.Any(), appL.Name).Return(&appL, nil)
	if !registration {
		ms.EXPECT().GetApplicationRole(gomock.Any(), appL.Name, username).Return(models.RoleAdmin, nil)
		ms.EXPECT().AddApplicationToUser(gomock.Any(), username, appL.Name, models.RoleAdmin)
	} else {
		ms.EXPECT().AddApplicationToUser(gomock.Any(), username, appL.Name, models.RoleUser)
	}

	// SUT
	err = au.AddUserToApplication(ctx, token, username, appL.Name, models.RoleAdmin)
	assert.NoError(t, err)
}

func TestAurum_AddUserToApplication(t *testing.T) {
	testAddToApplicationHelper(t, true)
	testAddToApplicationHelper(t, false)
}

func TestAurum_GetAccess(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"

	const app = "AnApp"
	appL := strings.ToLower(app)

	// Expect
	ms.EXPECT().GetApplicationRole(gomock.Any(), appL, username).Return(models.RoleAdmin, nil)

	// SUT
	au := Aurum{db: ms}
	resp, err := au.GetAccess(ctx, username, app)
	assert.NoError(t, err)

	assert.Equal(t, models.AccessStatus{
		ApplicationName: appL,
		Username:        username,
		AllowedAccess:   true,
		Role:            models.RoleAdmin,
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

	const app = "AnApp"
	appL := strings.ToLower(app)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetApplicationRole(gomock.Any(), appL, username).Return(models.RoleAdmin, nil)
	ms.EXPECT().SetApplicationRole(gomock.Any(), appL, target, models.RoleUser)

	// SUT
	err = au.SetAccess(ctx, token, app, target, models.RoleUser)
	assert.NoError(t, err)
}

func TestAurum_RemoveUserFromApplication(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	cfg := config.EphemeralConfig()

	ms := mock_store.NewMockAurumStore(ctrl)

	const username = "bob"
	const target = "wooloo"

	const app = "AnApp"
	appL := strings.ToLower(app)

	au := Aurum{db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey}

	token, err := jwt.GenerateJWT(username, false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetApplicationRole(gomock.Any(), appL, username).Return(models.RoleAdmin, nil)
	ms.EXPECT().RemoveApplicationFromUser(gomock.Any(), appL, target)

	// SUT
	err = au.RemoveUserFromApplication(ctx, token, app, target)
	assert.NoError(t, err)

}
