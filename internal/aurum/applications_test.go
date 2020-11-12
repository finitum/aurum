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

	au := Aurum{ db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey }

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

	au := Aurum{ db: ms, sk: cfg.SecretKey, pk: cfg.PublicKey }

	token, err := jwt.GenerateJWT("bob", false, cfg.SecretKey)
	assert.NoError(t, err)

	// Expect
	ms.EXPECT().GetApplicationRole(gomock.Any(), appL.Name, "bob").Return(models.RoleAdmin, nil)
	ms.EXPECT().RemoveApplication(gomock.Any(), appL.Name)

	// SUT
	err = au.RemoveApplication(ctx, token, appL.Name)
	assert.NoError(t, err)
}
