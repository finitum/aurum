package main

import (
	"context"
	"github.com/finitum/aurum/core/config"
	"github.com/finitum/aurum/pkg/store/dgraph"
	"github.com/finitum/aurum/services/aurum/routes"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {
	ctx := context.Background()

	dg, err := dgraph.New(ctx, "localhost:9080")
	if err != nil {
		log.Fatalf("Couldn't create dgraph client (%v)", err)
	}

	r := chi.NewRouter()

	cfg := config.GetConfig()

	rs := routes.NewRoutes(dg, cfg)

	r.Post("/signup", rs.SignUp)
	r.Post("/login", rs.Login)
	r.Get("/access/{app}/{user}", rs.Access)


}
