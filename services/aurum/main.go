package main

import (
	"context"
	"net/http"
	"time"

	"github.com/finitum/aurum/internal/aurum"
	"github.com/finitum/aurum/internal/cors"
	"github.com/finitum/aurum/pkg/config"
	"github.com/finitum/aurum/pkg/store/dgraph"
	"github.com/finitum/aurum/services/aurum/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {
	ctx := context.Background()
	cfg := config.GetConfig()

	log.Infof("Starting Aurum")

	var dg *dgraph.DGraph
	var err error
	for i := 0; i < 10; i++ {
		log.Infof("Connecting to DGraph")
		dg, err = dgraph.New(ctx, cfg.DgraphUrl)
		if err != nil {
			log.Errorf("Couldn't create Dgraph client, retrying in 3 seconds: %v", err)
			time.Sleep(3 * time.Second)
		} else {
			log.Infof("Connection with DGraph established")
			break
		}
	}
	if err != nil {
		log.Fatalf("Couldn't create Dgraph client: %v", err)
	}


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

	r.Get("/group/{group}/{user}", rs.GetAccess)

	r.Group(func(r chi.Router) {
		r.Use(rs.TokenExtractionMiddleware)

		r.Get("/user", rs.GetMe)
		r.Post("/user", rs.SetUser)
		r.Get("/user/{user}/groups", rs.GetGroupsForUser)

		// Group
		r.Post("/group", rs.AddGroup)
		r.Delete("/group/{group}", rs.RemoveGroup)

		r.Put("/group/{group}/{user}", rs.SetAccess)
		r.Post("/group/{group}/{user}", rs.AddUserToGroup)
		r.Delete("/group/{group}/{user}", rs.RemoveUserFromGroup)
	})

	srv := http.Server{
		Addr:         cfg.WebAddr,
		Handler:      r,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
