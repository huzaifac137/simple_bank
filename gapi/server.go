package gapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	"github.com/huzaifac137/simple_bank/pb"
	"github.com/huzaifac137/simple_bank/token"
	"github.com/huzaifac137/simple_bank/worker"
)

type Server struct {
	pb.SimpleBankServer
	symmetricKey    string
	tokenMaker      token.Maker
	store           db.Store
	router          *gin.Engine
	taskDistributer worker.TaskDistributor
}

func NewServer(symmetricKey string, store db.Store, taskDistributer worker.TaskDistributor) *Server {

	server := &Server{
		store:           store,
		symmetricKey:    symmetricKey,
		taskDistributer: taskDistributer,
	}

	maker, err := token.NewJWTMaker(symmetricKey)
	if err != nil {
		fmt.Printf("error creating New Maker , %s", err.Error())
	}

	server.tokenMaker = maker

	return server
}
