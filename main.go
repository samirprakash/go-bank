package main

import (
	"database/sql"
	"log"

	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"github.com/samirprakash/go-bank/api"
	db "github.com/samirprakash/go-bank/db/sqlc"
	"github.com/samirprakash/go-bank/util"
)

func main() {
	// load config from file or env vars using viper
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Not able to load config", err)
	}

	// connect to the database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	// create a database store for executing queries
	store := db.NewStore(conn)
	// create a server connected to the store
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server : ", err)
	}

	// start the server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
