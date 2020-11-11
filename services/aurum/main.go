package main

import (
	"context"
	"github.com/finitum/aurum/pkg/store"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {
	ctx := context.Background()

	dg, err := store.NewDGraph(ctx, "localhost:9080")
	if err != nil {
		log.Fatalf("Couldn't create dgraph client (%v)", err)
	}

	_ = dg
}
