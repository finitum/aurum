package dgraph

import (
	"context"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type DGraph struct {
	*dgo.Dgraph
}

func (dg DGraph) ClearAllImSure(ctx context.Context) error {
	return dg.Alter(ctx, &api.Operation{DropOp: api.Operation_ALL})
}

func New(ctx context.Context, address string) (*DGraph, error) {
	d, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	dg := dgo.NewDgraphClient(api.NewDgraphClient(d))

	if err := dg.Alter(ctx, &api.Operation{
		Schema: `
			type User {
				username
				password
				email
				groups
			}

			type Group {
				name
				allow_registration
			}

			username: string @index(hash) .
			password: string .
			email: string .
			groups: [uid] .

			name: string @index(hash) .
			allow_registration: bool .
		`,
	}); err != nil {
		return nil, errors.Wrap(err, "applying schema")
	}

	return &DGraph{dg}, nil
}
