package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/huzaifac137/simple_bank/api"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	gapi "github.com/huzaifac137/simple_bank/gapi"
	"github.com/huzaifac137/simple_bank/pb"
	"github.com/huzaifac137/simple_bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

//go:embed swagger/*
var swaggerFiles embed.FS

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	RunDBMigrations(config)

	store := db.NewStore(conn)

	go runGRPCGatewayServer(config, store)
	runGRPCServer(config, store)

}

func runGinServer(config *util.Config, store db.Store) {
	server := api.NewServer(config.SymmetricKey, store)
	e := server.Start(config.ServerAddress)
	if e != nil {
		log.Fatal().Err(e).Msg("Error starting server")
	}
}

func runGRPCServer(config *util.Config, store db.Store) {
	server := gapi.NewServer(config.SymmetricKey, store)

	unaryInterceptorOpt := gapi.GRPCUnaryInterceptor()
	grpcServer := grpc.NewServer(*unaryInterceptorOpt)

	pb.RegisterSimpleBankServer(grpcServer, server)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	log.Info().Str("port", config.ServerAddress).Msg("Starting gRPC server ")

	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatal().Err(err).Msg("grpc server failed")
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
		log.Fatal().Err(er).Msg("Failed to register grpc/http  gateway server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.FS(swaggerFiles))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	lis, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen using listener")
	}

	log.Info().Str("host:port", config.ServerAddress).Msg("Starting grpc/http gateway server ")

	handler := gapi.HTTPLogger(mux)
	err = http.Serve(lis, handler)

	if err != nil {
		log.Fatal().Err(err).Msg("grpc/http gateway server failed")
	}
}

func RunDBMigrations(config *util.Config) {
	migrateIns, err := migrate.New("file://db/migration", config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration instance")
	}

	err = migrateIns.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("cannot run migrations")
	}

	log.Info().Msg("Database migrations completed successfully")

}
