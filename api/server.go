package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	postRouter := router.Group("/api/accounts")
	postRouter.POST("/", server.createAccount)
	postRouter.GET("/:id", server.getAccount)
	postRouter.GET("/", server.listAccounts)
	server.router = router
	return server
}

func errorResponse(err error) *gin.H {
	return &gin.H{
		"error": err.Error(),
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
