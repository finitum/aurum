package store

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/finitum/aurum/pkg/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (dg DGraph) CreateUser(ctx context.Context, user *models.User) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	// query the database for the number of users that exist with either the same user id
	// or the same username
	query := `
		query q($uname: string) {
		  q(func: eq(username, $uname)) {
			count(uid)
		  }
		}
	`
	variables := map[string]string{
		"$uname": user.Username,
	}

	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return errors.Wrap(err, "query")
	}

	var r struct {
		Q []struct {
			Count int `json:"count"`
		} `json:"q"`
	}

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

	js, err := json.Marshal(dUser)
	if err != nil {
		return err
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson: js,
	}

	_, err = txn.Mutate(ctx, mu)
	return errors.Wrap(err, "mutate")
}

func (dg DGraph) getUser(ctx context.Context, user string) (*DGraphUser, error) {
	query := `
		query q($uname: string) {
		  q(func:eq(userID, $uname)) {
			uid
			username
			password
			email
		  }
		}
	`

	variables := map[string]string{"$uname": user}

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var r struct {
		Q []DGraphUser `json:"q"`
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	if len(r.Q) == 0 {
		return nil, errors.Errorf("user %s wasn't found", user)
	} else if len(r.Q) != 1 {
		return nil, errors.Errorf("expected one unique user %s, but found %d", user, len(r.Q))
	}

	return &r.Q[0], nil
}

func (dg DGraph) RemoveUser(ctx context.Context, username string) error {
	user, err := dg.getUser(ctx, username)
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

	_, err = dg.NewTxn().Mutate(ctx, mu)
	return errors.Wrap(err, "delete")
}

func (dg DGraph) GetUser(ctx context.Context, username string) (*models.User, error) {
	user, err := dg.getUser(ctx, username)
	if err != nil {
		return nil, err
	}

	return user.User, nil
}

func (dg DGraph) GetUsers(ctx context.Context) ([]models.User, error) {
	query := `
		{
			q(func: type(User)) {
				username
				password
				email
			}
		}
	`

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var r struct {
		Q []models.User `json:"q"`
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	return r.Q, nil
}

func (dg DGraph) AddUserToApplication(ctx context.Context, name string, appId uuid.UUID, role models.Role) error {
	panic("implement me")
}

func (dg DGraph) RemoveUserFromApplication(ctx context.Context, name string, appId uuid.UUID) error {
	panic("implement me")
}

func (dg DGraph) SetApplicationRole(ctx context.Context, name string, appId uuid.UUID, role models.Role) error {
	panic("implement me")
}

func (dg DGraph) GetApplicationRole(ctx context.Context, name string, appId uuid.UUID) (models.Role, error) {
	panic("implement me")
}

func (dg DGraph) CreateApplication(ctx context.Context, application *models.Application) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	// query the database for the number of applications that exist with either the same application id
	// or the same application name
	query := `
		query q($aid: string, $aname: string) {
		  q(func:type(Application)) @filter( eq(appID, $aid) OR eq(name, $aname)) {
				count(uid)
		  }
		}`

	variables := map[string]string{
		"$aid":   application.AppId.String(),
		"$aname": application.Name,
	}

	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return errors.Wrap(err, "query")
	}

	var r struct {
		Q []struct {
			Count int `json:"count"`
		} `json:"q"`
	}

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


	js, err := json.Marshal(dApplication)
	if err != nil {
		return err
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson: js,
	}

	_, err = txn.Mutate(ctx, mu)

	return errors.Wrap(err, "mutate")
}

func (dg DGraph) RemoveApplication(ctx context.Context, appId uuid.UUID) error {
	app, err := dg.getApplication(ctx, appId)
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

	_, err = dg.NewTxn().Mutate(ctx, mu)

	return errors.Wrap(err, "delete")
}

func (dg DGraph) getApplication(ctx context.Context, appId uuid.UUID) (*DGraphApplication, error) {
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
	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var r struct {
		Q []DGraphApplication `json:"q"`
	}

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

func (dg DGraph) GetApplication(ctx context.Context, appId uuid.UUID) (*models.Application, error) {
	app, err := dg.getApplication(ctx, appId)
	if err != nil {
		return nil, err
	}

	return &app.Application, nil
}

func (dg DGraph) GetApplications(ctx context.Context) ([]models.Application, error) {
	query := `
		{
			q(func: type(Application)) {
				appID
				name
			}
		}
	`

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var r struct {
		Q []models.Application `json:"q"`
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	return r.Q, nil
}
