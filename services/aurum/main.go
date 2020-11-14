package main

import (
	"context"
	"github.com/finitum/aurum/internal/cors"
	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/store/dgraph"
	"github.com/finitum/aurum/services/aurum/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {
	ctx := context.Background()

	dg, err := dgraph.New(ctx, "localhost:9080")
	if err != nil {
		log.Fatalf("Couldn't create Dgraph client: %v", err)
	}

	cfg := config.GetConfig()

	au, err := aurum.New(ctx, dg, cfg)
	if err != nil {
		log.Fatalf("Couldn't create Aurum client: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll)


	rs := routes.NewRoutes(au, cfg)

	r.Get("/pk", rs.PublicKey)

	r.Post("/signup", rs.SignUp)
	r.Post("/login", rs.Login)
	r.Post("/refresh", rs.Refresh)

	r.Get("/application/{app}/{user}", rs.GetAccess)

	r.Group(func(r chi.Router) {
		r.Use(rs.TokenExtractionMiddleware)

		r.Get("/user", rs.GetMe)
		r.Post("/user", rs.SetUser)

		// Application
		r.Post("/application", rs.AddApplication)
		r.Delete("/application", rs.RemoveApplication)

		r.Get("/application/{user}", rs.GetApplicationsForUser)

		r.Put("/application/{app}/{user}", rs.SetAccess)
		r.Post("/application/{app}/{user}", rs.AddUserToApplication)
		r.Delete("/application/{app}/{user}", rs.RemoveUserFromApplication)
	})

	log.Fatal(http.ListenAndServe(":8042", r))
}
