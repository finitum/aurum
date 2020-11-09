package main

import (
	"github.com/finitum/aurum/core/config"
	"github.com/finitum/aurum/core/db"
	"github.com/finitum/aurum/core/web"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Logrus has seven logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic.
	log.SetLevel(log.TraceLevel)
}

func main() {
	cfg := config.GetConfig()

	database := db.InitDB(db.InMemory)
	web.StartServer(cfg, database)
}
