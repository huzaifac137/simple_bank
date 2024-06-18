package main

import (
	"database/sql"
	"log"

	"github.com/huzaifac137/simple_bank/api"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	"github.com/huzaifac137/simple_bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	e := server.Start(config.ServerAddress)
	if e != nil {
		log.Fatal("cannot start server: ", e)
	}
}
