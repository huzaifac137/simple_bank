package main

import (
	"database/sql"
	"log"

	"github.com/huzaifac137/simple_bank/api"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:qwerty7867@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = ":8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	e := server.Start(serverAddress)
	if e != nil {
		log.Fatal("cannot start server: ", e)
	}
}
