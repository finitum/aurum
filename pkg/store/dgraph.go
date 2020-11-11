package store

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

type DGraph struct {
	*dgo.Dgraph
}

func (dg DGraph) ClearAllImSure(ctx context.Context) error {
	return dg.Alter(ctx, &api.Operation{DropOp: api.Operation_ALL})
}

func NewDGraph(ctx context.Context, address string) (*DGraph, error) {
	d, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	dg := dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)

	err = dg.Alter(ctx, &api.Operation{
		Schema: `
			type User {
				username
				password
				email
				applications
			}

			type Application {
				appID
				name
			}

			username: string @index(hash) .
			password: string .
			email: string .
			applications: [uid] .

			appID: string @index(hash) .
			name: string @index(hash) .
		`,
	})

	if err != nil {
		return nil, err
	}

	return &DGraph{dg }, nil
}
