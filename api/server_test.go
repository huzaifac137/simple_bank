package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	"github.com/huzaifac137/simple_bank/token"
)

func NewTestServer(store db.Store) *Server {

	server := &Server{
		store:        store,
		symmetricKey: "12345678901234567890123456789012",
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

	maker, err := token.NewJWTMaker("12345678901234567890123456789012")
	if err != nil {
		fmt.Printf("error creating New Maker , %s", err.Error())

	}

	server.tokenMaker = maker

	postRouter := router.Group("/api/accounts").Use(AuthMiddleware(maker))
	transferRouter := router.Group("/api/transfers").Use(AuthMiddleware(maker))
	userRouter := router.Group("/api/users")

	postRouter.POST("/", server.createAccount)
	postRouter.GET("/:id", server.getAccount)
	postRouter.GET("/", server.listAccounts)

	transferRouter.POST("/", server.CreateTransfer)

	userRouter.POST("/", server.createUser)
	userRouter.POST("/login", server.LoginUser)

	router.POST("/api/tokens/renew-access", server.RenewAccessToken)
	server.router = router
	return server
}
