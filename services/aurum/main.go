package main

import (
	"context"
	"github.com/finitum/aurum/internal/cors"
	"github.com/finitum/aurum/pkg/aurum"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/models"
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

	err = aurum.Initialize(ctx, dg)
	if err != nil {
		log.Fatalf("Couldn't initialize Aurum: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll)


	cfg := config.GetConfig()

	rs := routes.NewRoutes(dg, cfg)

	r.Get("/pk", rs.PublicKey)

	r.Post("/signup", rs.SignUp)
	r.Post("/login", rs.Login)
	r.Post("/refresh", rs.Refresh)

	r.Get("/access/{app}/{user}", rs.Access)

	r.Group(func(r chi.Router) {
		r.Use(rs.TokenVerificationMiddleware)

		r.Get("/user", rs.GetMe)
		r.Post("/user", rs.SetUser)

		r.Group(func(r chi.Router) {
			r.Use(rs.RoleMiddleware(models.RoleAdmin))

		})

	})

	log.Fatal(http.ListenAndServe(":8042", r))
}
