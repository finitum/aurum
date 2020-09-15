package main

import (
	"aurum/config"
	"aurum/db"
	"aurum/jwt/ecc"
	"aurum/web"
	"flag"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Logrus has seven logging levels: Trace, Debug, Info, Warning, Error, Fatal and Panic.
	log.SetLevel(log.TraceLevel)
}

func main() {
	generateKeys := flag.String("generate-keys", "none", "use to generate new keys. Options: [stdout, file, both]")

	flag.Parse()

	if *generateKeys == "none" {
		startServer()
	} else {
		ecfg := config.GetEnvConfig()
		ecc.KeyGenerationUtil(*generateKeys, ecfg.PublicKey, ecfg.SecretKeyPath)
	}
}

func startServer() {
	cfg := config.GetConfig()

	database := db.InitDB(db.INMEMORY)
	web.StartServer(cfg, database)
}
