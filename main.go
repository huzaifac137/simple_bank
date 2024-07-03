package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

	go runGRPCGatewayServer(config, store)
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

func runGRPCGatewayServer(config *util.Config, store db.Store) {
	server := gapi.NewServer(config.SymmetricKey, store)

	// create grpc muxxer to direct requests to grpc server
	grpcMux := runtime.NewServeMux()
	ctx, cncl := context.WithCancel(context.Background())
	defer cncl()
	er := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if er != nil {
		log.Fatalf("Failed to register gprc/http gateway handler: %v", er)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	lis, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting GRPC/HTTP GATEWAY server on", config.ServerAddress)

	err = http.Serve(lis, mux)

	if err != nil {
		log.Fatalf("grpc/http gateway server failed: %v", err)
	}
}
