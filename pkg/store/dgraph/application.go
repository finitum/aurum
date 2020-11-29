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

func (dg DGraph) getGroup(ctx context.Context, txn *dgo.Txn, name string) (*Group, error) {
	query := `
		query q($aname: string) {
		  q(func:eq(name, $aname)) {
			uid
			allow_registration
			name
		  }
		}
	`
	variables := map[string]string{"$aname": name}

	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var r struct {
		Q []Group `json:"q"`
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	if len(r.Q) == 0 {
		return nil, errors.Errorf("group %s wasn't found", name)
	} else if len(r.Q) != 1 {
		return nil, errors.Errorf("expected unique (one) group with name %s, but found %d", name, len(r.Q))
	}

	return &r.Q[0], nil
}

func (dg DGraph) GetGroup(ctx context.Context, name string) (*models.Group, error) {
	txn := dg.NewReadOnlyTxn().BestEffort()
	group, err := dg.getGroup(ctx, txn, name)
	if err != nil {
		return nil, err
	}

	return &group.Group, nil
}

func (dg DGraph) GetGroups(ctx context.Context) ([]models.Group, error) {
	query := `
		{
			q(func: type(Group)) {
				name
				allow_registration
			}
		}
	`

	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var r struct {
		Q []models.Group `json:"q"`
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	return r.Q, nil
}

func (dg DGraph) CreateGroup(ctx context.Context, group models.Group) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	// query the database for the number of groups that exist with either the same group id
	// or the same group name
	query := `
		query q($aname: string) {
		  q(func:eq(name, $aname)) {
			count(uid)
		  }
		}`

	variables := map[string]string{
		"$aname": group.Name,
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

	// If there exists 1 or more groups with this username, fail
	if len(r.Q) != 1 || r.Q[0].Count > 0 {
		return store.ErrExists
	}

	// Add the new group to the database
	dGroup := NewDGraphGroup(group)

	js, err := json.Marshal(dGroup)
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

func (dg DGraph) RemoveGroup(ctx context.Context, name string) error {
	txn := dg.NewTxn()

	group, err := dg.getGroup(ctx, txn, name)
	if err != nil {
		return errors.Wrap(err, "get user (internal)")
	}

	d := map[string]string{"uid": group.Uid}
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

func (dg DGraph) GetGroupsForUser(ctx context.Context, name string) ([]models.GroupWithRole, error) {
	query := `
query q($uname: string) {
  q(func: type(User)) @filter(eq(username, $uname)) {
	username
   	groups @facets(role:role) {
      name
	  allow_registration
  	} 
  }
}`

	variables := map[string]string{
		"$uname": name,
	}
	txn := dg.NewReadOnlyTxn().BestEffort()
	resp, err := txn.QueryWithVars(ctx, query, variables)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}

	var r struct {
		Q []struct {
			Groups []models.GroupWithRole `json:"groups"`
		} `json:"q"`
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	} else if len(r.Q) != 1 {
		return nil, errors.Wrap(err, "how the hell did this happen???")
	} else if len(r.Q[0].Groups) == 0 {
		return nil, store.ErrNotExists
	}

	return r.Q[0].Groups, nil
}
