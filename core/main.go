package main

import (
	"aurum/config"
	"aurum/db"
	"aurum/web"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Logrus has seven logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic.
	log.SetLevel(log.TraceLevel)
}

func main() {
	cfgbuilder := config.Builder{}
	cfg := cfgbuilder.SetDefault().SetFromEnvironment().FindKeys(false).Build()

	database := db.InitDB(db.INMEMORY)
	web.StartServer(cfg, database)
}
