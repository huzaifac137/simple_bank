package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/huzaifac137/simple_bank/api"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	gapi "github.com/huzaifac137/simple_bank/gapi"
	"github.com/huzaifac137/simple_bank/pb"
	"github.com/huzaifac137/simple_bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
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

	//runGinServer(config, store)
	runGRPCServer(config, store)

}

func runGinServer(config *util.Config, store db.Store) {
	server := api.NewServer(config.SymmetricKey, store)
	e := server.Start(config.ServerAddress)
	if e != nil {
		log.Fatal("cannot start server: ", e)
	}
}

func runGRPCServer(config *util.Config, store db.Store) {
	server := gapi.NewServer(config.SymmetricKey, store)

	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting gRPC server on", config.GRPCServerAddress)

	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
