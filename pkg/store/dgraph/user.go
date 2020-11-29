package dgraph

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store"
	"github.com/pkg/errors"
)

func (dg DGraph) getUser(ctx context.Context, txn *dgo.Txn, user string) (User, error) {
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

	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return User{}, errors.Wrap(err, "query")
	}

	var r struct {
		Q []User `json:"q"`
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return User{}, errors.Wrap(err, "json unmarshal")
	}

	if len(r.Q) == 0 {
		return User{}, errors.Errorf("user %s wasn't found", user)
	} else if len(r.Q) != 1 {
		return User{}, errors.Errorf("expected one unique user %s, but found %d", user, len(r.Q))
	}

	return r.Q[0], nil
}

func (dg DGraph) getUserWithGroups(ctx context.Context, txn *dgo.Txn, user string, name string) (User, error) {
	query := `
query q($uname: string, $aname: string) {
	User(func:eq(username, $uname)) {
		uid
		groups @facets @filter(eq(name, $aname)) {
			uid
		}
	}
}
`
	variables := map[string]string{"$uname": user, "$aname": name}

	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return User{}, errors.Wrap(err, "query")
	}

	var r struct {
		User []User
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return User{}, errors.Wrap(err, "json unmarshal")
	}

	if len(r.User) != 1 || len(r.User[0].Groups) != 1 {
		return User{}, store.ErrNotExists
	}

	return r.User[0], nil
}

func (dg DGraph) GetUser(ctx context.Context, name string) (models.User, error) {
	txn := dg.NewReadOnlyTxn()

	user, err := dg.getUser(ctx, txn, name)
	if err != nil {
		return models.User{}, err
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

func (dg DGraph) CreateUser(ctx context.Context, user models.User) error {
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
		return store.ErrExists
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

func (dg DGraph) SetUser(ctx context.Context, user models.User) (models.User, error) {
	txn := dg.NewTxn()

	currUser, err := dg.getUser(ctx, txn, user.Username)
	if err != nil {
		return models.User{}, err
	}

	if user.Password != "" {
		currUser.Password = user.Password
	}

	if user.Email != "" {
		currUser.Email = user.Email
	}

	js, err := json.Marshal(&currUser)
	if err != nil {
		return models.User{}, err
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		SetJson:   js,
		CommitNow: true,
	})
	if err != nil {
		return models.User{}, err
	}

	return currUser.User, nil
}

func (dg DGraph) RemoveUser(ctx context.Context, username string) error {
	txn := dg.NewTxn()

	user, err := dg.getUser(ctx, txn, username)
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

	_, err = txn.Mutate(ctx, mu)
	return errors.Wrap(err, "delete")
}

func (dg DGraph) GetGroupRole(ctx context.Context, group string, user string) (models.Role, error) {
	txn := dg.NewReadOnlyTxn().BestEffort()

	u, err := dg.getUserWithGroups(ctx, txn, user, group)
	if err != nil {
		return 0, err
	}

	return u.Groups[0].Role, nil
}

func (dg DGraph) AddGroupToUser(ctx context.Context, user string, group string, role models.Role) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	q := `
query q($uname: string, $gname: string) {
  User(func:eq(username, $uname)) {
    uid
  }

  Group(func:eq(name, $gname)) {
  	uid
  }
}
`
	var r struct {
		User  []User
		Group []Group
	}

	vars := map[string]string{
		"$uname": user,
		"$gname": group,
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Json, &r); err != nil {
		return err
	}

	if len(r.User) != 1 || len(r.Group) != 1 {
		return errors.New("Couldn't find user or group")
	}

	r.Group[0].Role = role
	r.User[0].Groups = []Group{r.Group[0]}

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

func (dg DGraph) SetGroupRole(ctx context.Context, group string, user string, role models.Role) error {
	return dg.AddGroupToUser(ctx, user, group, role)
}

func (dg DGraph) RemoveGroupFromUser(ctx context.Context, group string, user string) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	u, err := dg.getUserWithGroups(ctx, txn, user, group)

	js, err := json.Marshal(&u)
	if err != nil {
		return err
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		DeleteJson: js,
		CommitNow:  true,
	})

	return errors.Wrap(err, "mutate")
}

func (dg DGraph) CountUsers(ctx context.Context) (int, error) {
	query := `
{
	Q(func: type(User)) {
		count(uid)
	}
}
	`

	txn := dg.NewReadOnlyTxn().BestEffort()
	defer txn.Discard(ctx)

	resp, err := txn.Query(ctx, query)
	if err != nil {
		return -1, errors.Wrap(err, "query")
	}

	var r struct {
		Q []struct {
			Count int `json:"count"`
		}
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return -1, errors.Wrap(err, "json unmarshal")
	}

	if len(r.Q) != 1 {
		return -1, errors.New("unexpected multiple results in count")
	}

	return r.Q[0].Count, nil
}
