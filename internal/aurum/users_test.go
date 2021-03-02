package aurum

import (
	"context"
	"github.com/finitum/aurum/pkg/store"
	"reflect"
	"testing"

	"github.com/finitum/aurum/internal/hash"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store/mock_store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAurum_SignUp(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	u := models.User{
		Username: "user",
		Password: "wH6VLfolKTUb",
		Email:    "email",
	}

	ctxT := reflect.TypeOf(ctx)

	ms.EXPECT().CreateUser(gomock.AssignableToTypeOf(ctxT), gomock.Any()).Do(func(_ context.Context, gu models.User) {
		assert.True(t, hash.CheckPasswordHash(u.Password, gu.Password))
	}).Return(nil)
	ms.EXPECT().AddGroupToUser(gomock.AssignableToTypeOf(ctxT), u.Username, AurumName, models.RoleUser)

	au := Aurum{db: ms}
	// SUT
	err := au.SignUp(ctx, u)
	assert.NoError(t, err)
}

func TestAurum_SignUp_EExists(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()
	ms := mock_store.NewMockAurumStore(ctrl)

	u := models.User{
		Username: "user",
		Password: "wH6VLfolKTUb",
		Email:    "email",
	}

	ms.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(store.ErrExists)

	au := Aurum{db: ms}
	// SUT
	err := au.SignUp(ctx, u)
	assert.Equal(t, store.ErrExists, err)
}

func TestAurum_SignUp_WeakPass(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()
	ms := mock_store.NewMockAurumStore(ctrl)

	u := models.User{
		Username: "user",
		Password: "123",
		Email:    "email",
	}

	au := Aurum{db: ms}
	// SUT
	err := au.SignUp(ctx, u)
	assert.Equal(t, ErrWeakPassword, err)
}


func TestAurum_Login(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	cfg := config.EphemeralConfig()

	au := Aurum{db: ms, pk: cfg.PublicKey, sk: cfg.SecretKey}

	u := models.User{
		Username: "user",
		Password: "wH6VLfolKTUb",
		Email:    "email",
	}

	var err error
	hu := u
	hu.Password, err = hash.HashPassword(u.Password)
	assert.NoError(t, err)

	ms.EXPECT().GetUser(gomock.Any(), u.Username).Return(hu, nil)

	// SUT
	tp, err := au.Login(ctx, u)
	assert.NoError(t, err)

	lt, err := jwt.VerifyJWT(tp.LoginToken, cfg.PublicKey)
	assert.NoError(t, err)
	assert.False(t, lt.Refresh)
	assert.Equal(t, u.Username, lt.Username)

	rt, err := jwt.VerifyJWT(tp.RefreshToken, cfg.PublicKey)
	assert.NoError(t, err)
	assert.True(t, rt.Refresh)
}

func TestAurum_Login_WrongPass(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	cfg := config.EphemeralConfig()

	au := Aurum{db: ms, pk: cfg.PublicKey, sk: cfg.SecretKey}

	u := models.User{
		Username: "user",
		Password: "wH6VLfolKTUb",
		Email:    "email",
	}

	var err error
	hu := u
	hu.Password, err = hash.HashPassword("fakepass")
	assert.NoError(t, err)

	ms.EXPECT().GetUser(gomock.Any(), u.Username).Return(hu, nil)

	// SUT
	tp, err := au.Login(ctx, u)
	assert.Empty(t, tp)
	assert.Error(t, err)
}

func TestAurum_RefreshToken(t *testing.T) {
	cfg := config.EphemeralConfig()

	au := Aurum{pk: cfg.PublicKey, sk: cfg.SecretKey}

	// SUT
	tp, err := jwt.GenerateJWTPair("jeff", cfg.SecretKey)
	assert.NoError(t, err)

	old := tp
	err = au.RefreshToken(&tp)
	assert.NoError(t, err)

	assert.Equal(t, old.RefreshToken, tp.RefreshToken)
	assert.NotEqual(t, old.LoginToken, tp.LoginToken)

	lt, err := jwt.VerifyJWT(tp.LoginToken, cfg.PublicKey)
	assert.NoError(t, err)
	assert.False(t, lt.Refresh)
	assert.Equal(t, "jeff", lt.Username)
}

func TestAurum_RefreshInvalidToken(t *testing.T) {
	cfg := config.EphemeralConfig()

	au := Aurum{pk: cfg.PublicKey, sk: cfg.SecretKey}

	// SUT
	tp := jwt.TokenPair{
		LoginToken:   "invalid",
		RefreshToken: "invalid",
	}

	err := au.RefreshToken(&tp)
	assert.Error(t, err)
}


func TestAurum_GetUser(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	cfg := config.EphemeralConfig()

	au := Aurum{db: ms, pk: cfg.PublicKey, sk: cfg.SecretKey}

	u := models.User{
		Username: "user",
		Password: "wH6VLfolKTUb",
		Email:    "email",
	}

	ms.EXPECT().GetUser(gomock.Any(), u.Username).Return(u, nil)

	tp, err := jwt.GenerateJWTPair(u.Username, cfg.SecretKey)
	assert.NoError(t, err)
	// SUT
	gu, err := au.GetUser(ctx, tp.LoginToken)
	assert.NoError(t, err)

	assert.Equal(t, models.User{
		Username: u.Username,
		Email:    u.Email,
	}, gu)
}

func TestAurum_UpdateUser(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	cfg := config.EphemeralConfig()

	au := Aurum{db: ms, pk: cfg.PublicKey, sk: cfg.SecretKey}

	u := models.User{
		Username: "user",
		Password: "wH6VLfolKTUb",
		Email:    "email",
	}

	ms.EXPECT().SetUser(gomock.Any(), gomock.Any()).Do(func(_ context.Context, user models.User) {
		assert.True(t, hash.CheckPasswordHash(u.Password, user.Password))
	}).Return(u, nil)

	tp, err := jwt.GenerateJWTPair(u.Username, cfg.SecretKey)
	assert.NoError(t, err)
	// SUT
	gu, err := au.UpdateUser(ctx, tp.LoginToken, u)

	assert.Equal(t, models.User{
		Username: u.Username,
		Email:    u.Email,
	}, gu)
}


func TestAurum_UpdateUser_WeakPassword(t *testing.T) {
	ctx := context.Background()
	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	ms := mock_store.NewMockAurumStore(ctrl)

	cfg := config.EphemeralConfig()

	au := Aurum{db: ms, pk: cfg.PublicKey, sk: cfg.SecretKey}

	u := models.User{
		Username: "user",
		Password: "weak",
	}

	tp, err := jwt.GenerateJWTPair(u.Username, cfg.SecretKey)
	assert.NoError(t, err)
	// SUT
	gu, err := au.UpdateUser(ctx, tp.LoginToken, u)
	assert.Equal(t, ErrWeakPassword, err)
	assert.Empty(t, gu)
}
