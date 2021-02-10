package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/samirprakash/go-bank/api"
	db "github.com/samirprakash/go-bank/db/sqlc"
	"github.com/samirprakash/go-bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Not able to load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
