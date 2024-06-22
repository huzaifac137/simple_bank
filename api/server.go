package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
)

type Server struct {
	symmetricKey string
	store        db.Store
	router       *gin.Engine
}

func NewServer(symmetricKey string, store db.Store) *Server {

	server := &Server{
		store:        store,
		symmetricKey: symmetricKey,
	}

	router := gin.Default()

	// error validator instance
	validate := validator.New()

	// register validator with currency tag
	validate.RegisterValidation("currency", customCurrencyValidator)

	// Register the validator with Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", customCurrencyValidator)
	}

	postRouter := router.Group("/api/accounts")
	transferRouter := router.Group("/api/transfers")
	userRouter := router.Group("/api/users")

	postRouter.POST("/", server.createAccount)
	postRouter.GET("/:id", server.getAccount)
	postRouter.GET("/", server.listAccounts)

	transferRouter.POST("/", server.CreateTransfer)

	userRouter.POST("/", server.createUser)
	userRouter.POST("/login", server.LoginUser)
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
