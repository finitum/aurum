package main

import (
	"aurum/config"
	"aurum/db"
	"aurum/web"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	cfg := new(config.Builder).SetDefault().Build()
	database := db.InitDB("inmemory")
	web.StartServer(cfg, database)
}
