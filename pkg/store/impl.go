package store

import (
	"encoding/json"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/finitum/aurum/pkg/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (dg DGraph) CreateUser(user *models.User) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(dg.ctx)

	// query the database for the number of users that exist with either the same user id
	// or the same username
	query := `
		query q($uid: string, $uname: string) {
		  q(func:type(User)) @filter( eq(userID, $uid) OR eq(username, $uname)) {
				count(uid)
		  }
		}
	`
	variables := map[string]string{
		"$uid":   user.UserId.String(),
		"$uname": user.Username,
	}

	resp, err := txn.QueryWithVars(dg.ctx, query, variables)
	if err != nil {
		return errors.Wrap(err, "query")
	}

	type Root struct {
		Q []struct {
			Count int `json:"count"`
		} `json:"q"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	// If there exists 1 or more users with this username, fail
	if len(r.Q) != 1 || r.Q[0].Count > 0 {
		return errors.Errorf("user %s exists", user.Username)
	}

	// Add the new user to the database
	dUser := NewDGraphUser(user)

	mu := &api.Mutation{
		CommitNow: true,
	}

	js, err := json.Marshal(dUser)
	if err != nil {
		return err
	}

	mu.SetJson = js

	_, err = txn.Mutate(dg.ctx, mu)
	if err != nil {
		return errors.Wrap(err, "mutate")
	}

	return nil
}

func (dg DGraph) RemoveUser(userId uuid.UUID) error {
	user, err := dg.getUserInternal(userId)
	if err != nil {
		return errors.Wrap(err, "get user (internal)")
	}

	d := map[string]string{"uid": user.Uid}
	js, err := json.Marshal(d)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	mu := &api.Mutation{
		CommitNow:  true,
		DeleteJson: js,
	}

	_, err = dg.NewTxn().Mutate(dg.ctx, mu)

	return errors.Wrap(err, "delete")
}

func (dg DGraph) getUserInternal(userID uuid.UUID) (*DGraphUser, error) {
	query := `
		query q($uid: string) {
		  q(func:eq(userID, $uid)) {
			uid
			userID
			username
			password
			email
		  }
		}
	`

	variables := map[string]string{"$uid": userID.String()}

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.QueryWithVars(dg.ctx, query, variables)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	type Root struct {
		Q []DGraphUser `json:"q"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	if len(r.Q) == 0 {
		return nil, errors.Errorf("user with user id %s wasn't found", userID)
	} else if len(r.Q) != 1 {
		return nil, errors.Errorf("expected unique (one) user id %s, but found %d", userID, len(r.Q))
	}

	return &r.Q[0], nil
}

func (dg DGraph) GetUser(userId uuid.UUID) (*models.User, error) {
	user, err := dg.getUserInternal(userId)
	if err != nil {
		return nil, err
	}

	return user.User, nil
}

func (dg DGraph) GetUsers() ([]models.User, error) {
	query := `
		{
			q(func: type(User)) {
				userID
				username
				password
				email
			}
		}
	`

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.Query(dg.ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	type Root struct {
		Q []models.User `json:"q"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	return r.Q, nil
}

func (dg DGraph) AddUserToApplication(userId uuid.UUID, appId uuid.UUID, role models.Role) error {
	panic("implement me")
}

func (dg DGraph) RemoveUserFromApplication(userId uuid.UUID, appId uuid.UUID) error {
	panic("implement me")
}

func (dg DGraph) SetApplicationRole(userId uuid.UUID, appId uuid.UUID, role models.Role) error {
	panic("implement me")
}

func (dg DGraph) GetApplicationRole(userId uuid.UUID, appId uuid.UUID) (models.Role, error) {
	panic("implement me")
}

func (dg DGraph) CreateApplication(application *models.Application) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(dg.ctx)

	// query the database for the number of users that exist with either the same user id
	// or the same username
	query := `
		query q($aid: string, $aname: string) {
		  q(func:type(Application)) @filter( eq(appID, $aid) OR eq(name, $aname)) {
				count(uid)
		  }
		}
	`
	variables := map[string]string{
		"$aid":   application.AppId.String(),
		"$aname": application.Name,
	}

	resp, err := txn.QueryWithVars(dg.ctx, query, variables)
	if err != nil {
		return errors.Wrap(err, "query")
	}

	type Root struct {
		Q []struct {
			Count int `json:"count"`
		} `json:"q"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	// If there exists 1 or more users with this username, fail
	if len(r.Q) != 1 || r.Q[0].Count > 0 {
		return errors.Errorf("application %s exists", application.Name)
	}

	// Add the new user to the database
	dApplication := NewDGraphApplication(application)

	mu := &api.Mutation{
		CommitNow: true,
	}

	js, err := json.Marshal(dApplication)
	if err != nil {
		return err
	}

	mu.SetJson = js

	_, err = txn.Mutate(dg.ctx, mu)
	if err != nil {
		return errors.Wrap(err, "mutate")
	}

	return nil
}

func (dg DGraph) RemoveApplication(appId uuid.UUID) error {
	app, err := dg.getApplicationInternal(appId)
	if err != nil {
		return errors.Wrap(err, "get user (internal)")
	}

	d := map[string]string{"uid": app.Uid}
	js, err := json.Marshal(d)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	mu := &api.Mutation{
		CommitNow:  true,
		DeleteJson: js,
	}

	_, err = dg.NewTxn().Mutate(dg.ctx, mu)

	return errors.Wrap(err, "delete")
}

func (dg DGraph) getApplicationInternal(appId uuid.UUID) (*DGraphApplication, error) {
	query := `
		query q($aid: string) {
		  q(func:eq(appID, $aid)) {
			uid
			appID
			name
		  }
		}
	`

	variables := map[string]string{"$aid": appId.String()}

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.QueryWithVars(dg.ctx, query, variables)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	type Root struct {
		Q []DGraphApplication `json:"q"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	if len(r.Q) == 0 {
		return nil, errors.Errorf("application with app id %s wasn't found", appId)
	} else if len(r.Q) != 1 {
		return nil, errors.Errorf("expected unique (one) application id %s, but found %d", appId, len(r.Q))
	}

	return &r.Q[0], nil
}

func (dg DGraph) GetApplication(appId uuid.UUID) (*models.Application, error) {
	app, err := dg.getApplicationInternal(appId)
	if err != nil {
		return nil, err
	}

	return &app.Application, nil
}

func (dg DGraph) GetApplications() ([]models.Application, error) {
	query := `
		{
			q(func: type(Application)) {
				appID
				name
			}
		}
	`

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.Query(dg.ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	type Root struct {
		Q []models.Application `json:"q"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	return r.Q, nil
}
