package main

import (
	"context"
	"github.com/finitum/aurum/pkg/models"
	"github.com/finitum/aurum/pkg/store/dgraph"
	"github.com/google/uuid"
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

	user := models.User{
		Username: "user",
		Password: "pass",
		Email:    "email",
	}
	if err := dg.CreateUser(ctx, &user); err != nil {
		log.Fatal(err)
	}

	app := models.Application{
		AppId: uuid.New(),
		Name:  "aurum",
	}
	if err := dg.CreateApplication(ctx, &app); err != nil {
		log.Fatal(err)
	}

	if err := dg.AddUserToApplication(ctx, user.Username, app.AppId, models.RoleAdmin); err != nil {
		log.Fatal(err)
	}
}
