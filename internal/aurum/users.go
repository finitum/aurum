package aurum

import (
	"context"
	"github.com/finitum/aurum/internal/hash"
	"github.com/finitum/aurum/internal/passwords"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/pkg/errors"
)

func (au Aurum) SignUp(ctx context.Context, user models.User) error {
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

	if err := au.db.CreateUser(ctx, user); err != nil {
		if err == store.ErrExists {
			return err
		}

		return errors.Wrap(err, "failed creating user in database")
	}

	if err := au.db.AddGroupToUser(ctx, user.Username, AurumName, models.RoleUser); err != nil {
		return errors.Wrap(err, "couldn't add user to aurum group")
	}

	return nil
}

func (au Aurum) Login(ctx context.Context, user models.User) (jwt.TokenPair, error) {
	dbu, err := au.db.GetUser(ctx, user.Username)
	if err != nil {
		return jwt.TokenPair{}, errors.Wrap(err, "getting user from db failed")
	}

	if !hash.CheckPasswordHash(user.Password, dbu.Password) {
		return jwt.TokenPair{}, errors.New("invalid password")
	}

	return jwt.GenerateJWTPair(dbu.Username, au.sk)
}

func (au Aurum) RefreshToken(tp *jwt.TokenPair) error {
	if tp.RefreshToken == "" {
		return ErrInvalidInput
	}

	claims, err := jwt.VerifyJWT(tp.RefreshToken, au.pk)
	if err != nil {
		return errors.Wrap(err, "verification error")
	}

	newtoken, err := jwt.GenerateJWT(claims.Username, false, au.sk)
	if err != nil {
		return errors.Wrap(err, "jwt generation error")
	}

	tp.LoginToken = newtoken

	return nil
}

func (au Aurum) GetUser(ctx context.Context, token string) (models.User, error) {
	claims, err := au.checkToken(token)
	if err != nil {
		return models.User{}, err
	}

	user, err := au.db.GetUser(ctx, claims.Username)
	if err != nil {
		return models.User{}, err
	}

	user.Password = ""

	return user, nil
}

func (au Aurum) UpdateUser(ctx context.Context, token string, user models.User) (models.User, error) {

	claims, err := au.checkToken(token)
	if err != nil {
		return models.User{}, err
	}

	user.Username = claims.Username

	if user.Password != "" {
		if !passwords.CheckStrength(user.Password, []string{user.Username, user.Email}) {
			return models.User{}, ErrWeakPassword
		}

		hashed, err := hash.HashPassword(user.Password)
		if err != nil {
			return models.User{}, err
		}

		user.Password = hashed
	}

	user, err = au.db.SetUser(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	user.Password = ""

	return user, nil
}
