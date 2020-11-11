package aurum

import (
	"context"
	"github.com/finitum/aurum/internal/hash"
	"github.com/finitum/aurum/internal/passwords"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrInvalidInput = errors.New("password is too weak")
var ErrWeakPassword = errors.New("password is too weak")

func SignUp(ctx context.Context, db store.AurumStore, user models.User) error {
	if user.Username == "" {
		return ErrInvalidInput
	}

	if !passwords.CheckStrength(user.Password, []string{user.Username, user.Email}) {
		return ErrWeakPassword
	}

	hashed, err := hash.HashPassword(user.Password)
	if err != nil {
		return errors.Wrap(err, "hashing failed")
	}

	user.Password = hashed

	if err := db.CreateUser(ctx, &user); err != nil {
		if err == store.ErrExists {
			return err
		}

		return errors.Wrap(err, "failed creating user in database")
	}

	return nil
}

func Login(ctx context.Context, db store.AurumStore, user models.User, key ecc.SecretKey) (jwt.TokenPair, error) {
	dbu, err := db.GetUser(ctx, user.Username)
	if err != nil {
		return jwt.TokenPair{}, errors.Wrap(err, "getting user from db failed")
	}

	if !hash.CheckPasswordHash(user.Password, dbu.Password) {
		return jwt.TokenPair{}, errors.New("invalid password")
	}

	return jwt.GenerateJWTPair(dbu.Username, key)
}

func Access(ctx context.Context, db store.AurumStore, user string, appid uuid.UUID) (models.AccessResponse, error) {
	role, err := db.GetApplicationRole(ctx, user, appid);

	if err == store.ErrNotExists {
		return models.AccessResponse{
			ApplicationID: appid,
			Username:      user,
			AllowedAccess: false,
		}, nil

	} else if err != nil {
		return models.AccessResponse{}, err
	}

	return models.AccessResponse{
		ApplicationID: appid,
		Username:      user,
		AllowedAccess: true,
		Role:          role,
	}, nil
}

func RefreshToken(tp *jwt.TokenPair, pk ecc.PublicKey, sk ecc.SecretKey) error {
	if tp.RefreshToken == "" {
		return ErrInvalidInput
	}
	
	claims, err := jwt.VerifyJWT(tp.RefreshToken, pk)
	if err != nil {
		return errors.Wrap(err, "verification error")
	}
	
	newtoken, err := jwt.GenerateJWT(claims.Username, false, sk)
	if err != nil {
		return errors.Wrap(err, "jwt generation error")
	}

	tp.LoginToken = newtoken

	return nil
}

func GetUser(ctx context.Context, db store.AurumStore, user string) (models.User, error) {
	return db.GetUser(ctx, user)
}

func UpdateUser(ctx context.Context, db store.AurumStore, user models.User) error {

	if !passwords.CheckStrength(user.Password, []string{user.Username, user.Email}) {
		return ErrWeakPassword
	}

	hashed, err := hash.HashPassword(user.Password)
	if err != nil {
		return errors.Wrap(err, "hashing failed")
	}

	user.Password = hashed

	return db.SetUser(ctx, &user)
}
