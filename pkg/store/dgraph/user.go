package dgraph

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/finitum/aurum/core/db"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (dg DGraph) getUser(ctx context.Context, user string) (*User, error) {
	query := `
query q($uname: string) {
	q(func:eq(username, $uname)) {
		uid
		username
		password
		email
	}
}`

	variables := map[string]string{"$uname": user}

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var r struct {
		Q []User `json:"q"`
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

func (dg DGraph) CreateUser(ctx context.Context, user *models.User) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	// query the database for the number of users that exist with either the same user id
	// or the same username
	query := `
		query q($uname: string) {
		  Q(func: eq(username, $uname)) {
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
		}
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}

	// If there exists 1 or more users with this username, fail
	if len(r.Q) != 1 || r.Q[0].Count > 0 {
		return db.ErrExists
	}

	// Add the new user to the database
	dUser := NewDGraphUser(user)

	js, err := json.Marshal(dUser)
	if err != nil {
		return err
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   js,
	}

	_, err = txn.Mutate(ctx, mu)
	return errors.Wrap(err, "mutate")
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

func (dg DGraph) GetApplicationRole(ctx context.Context, name string, appId uuid.UUID) (models.Role, error) {
	query := `
query q($uname: string, $aid: string) {
	User(func:eq(username, $uname)) {
		uid
		applications @facets @filter(eq(appID, $aid)) {
			uid
		}
	}
}
`
	variables := map[string]string{"$uname": name, "$aid": appId.String()}

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return 0, errors.Wrap(err, "query")
	}

	var r struct {
		User []User
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return 0, errors.Wrap(err, "json unmarshal")
	}

	if len(r.User) != 1 || len(r.User[0].Applications) != 1 {
		return 0, store.ErrNotExists
	}

	return r.User[0].Applications[0].Role, nil
}

func (dg DGraph) SetApplicationRole(ctx context.Context, name string, appId uuid.UUID, role models.Role) error {
	return dg.AddApplicationToUser(ctx, name, appId, role)
}


func (dg DGraph) AddApplicationToUser(ctx context.Context, name string, appId uuid.UUID, role models.Role) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	q := `
query q($uname: string, $aid: string) {
  User(func:eq(username, $uname)) {
    uid
  }

  App(func:eq(appID, $aid)) {
  	uid
	appID
  }
}
`
	var r struct {
		User []User
		App  []Application
	}

	vars := map[string]string{
		"$uname": name,
		"$aid":   appId.String(),
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Json, &r); err != nil {
		return err
	}

	if len(r.User) != 1 || len(r.App) != 1 {
		return errors.New("Couldn't find user or application")
	}

	r.App[0].Role = role
	r.User[0].Applications = []Application{ r.App[0] }

	js, err := json.Marshal(&r.User[0])
	if err != nil {
		return nil
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   js,
	}

	if _, err := txn.Mutate(ctx, mu); err != nil {
		return err
	}

	return nil
}

func (dg DGraph) RemoveApplicationFromUser(ctx context.Context, name string, appId uuid.UUID) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	q := `
query q($uname: string, $aid: string) {
  User(func:eq(username, $uname)) {
    uid
	applications @filter(eq(appID, $aid)) {
		uid
	}
  }
}
`

	vars := map[string]string{
		"$uname": name,
		"$aid":   appId.String(),
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return err
	}

	var r struct {
		User []User
	}

	if err := json.Unmarshal(resp.Json, &r); err != nil {
		return err
	}

	if len(r.User) != 1 || len(r.User[0].Applications) != 1 {
		return errors.New("Couldn't find user or application")
	}

	js, err := json.Marshal(r.User[0])
	if err != nil {
		return err
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		DeleteJson:           js,
		CommitNow:            true,
	})

	return errors.Wrap(err, "mutate")
}


