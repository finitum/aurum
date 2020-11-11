package dgraph

import (
	"context"
	"encoding/json"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/finitum/aurum/pkg/models"
	"github.com/pkg/errors"
)

func (dg DGraph) getApplication(ctx context.Context, txn *dgo.Txn, name string) (*Application, error) {
	query := `
		query q($aname: string) {
		  q(func:eq(name, $aname)) {
			uid
			appID
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
		Q []Application `json:"q"`
	}

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	if len(r.Q) == 0 {
		return nil, errors.Errorf("application %s wasn't found", name)
	} else if len(r.Q) != 1 {
		return nil, errors.Errorf("expected unique (one) application with name %s, but found %d", name, len(r.Q))
	}

	return &r.Q[0], nil
}

func (dg DGraph) GetApplication(ctx context.Context, name string) (*models.Application, error) {
	txn := dg.NewReadOnlyTxn().BestEffort()
	app, err := dg.getApplication(ctx, txn, name)
	if err != nil {
		return nil, err
	}

	return &app.Application, nil
}

func (dg DGraph) GetApplications(ctx context.Context) ([]models.Application, error) {
	query := `
		{
			q(func: type(Application)) {
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

func (dg DGraph) CreateApplication(ctx context.Context, application models.Application) error {
	// start a new transaction
	txn := dg.NewTxn()
	defer txn.Discard(ctx)

	// query the database for the number of applications that exist with either the same application id
	// or the same application name
	query := `
		query q($aname: string) {
		  q(func:eq(name, $aname)) {
			count(uid)
		  }
		}`

	variables := map[string]string{
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
		SetJson:   js,
	}

	_, err = txn.Mutate(ctx, mu)

	return errors.Wrap(err, "mutate")
}

func (dg DGraph) RemoveApplication(ctx context.Context, name string) error {
	txn := dg.NewTxn()

	app, err := dg.getApplication(ctx, txn, name)
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

	_, err = txn.Mutate(ctx, mu)

	return errors.Wrap(err, "delete")
}
