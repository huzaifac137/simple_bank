package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/huzaifac137/simple_bank/token"
)

const (
	authorizationHeaderKey     = "authorization"
	authorizationHeaderType    = "Bearer"
	authorizationHeaderPayload = "authorization_payload"
)

func AuthMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(authorizationHeaderKey)

		if len(authorization) == 0 {
			err := errors.New("authorization token is required")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}

		splitHeader := strings.Fields(authorization)

		headerType := splitHeader[0]

		if headerType != authorizationHeaderType {
			err := fmt.Errorf("Unspported authorization header type provided : %s", headerType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}

		tokenString := splitHeader[1]
		payload, err := maker.VerifyToken(tokenString)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}

		ctx.Set(authorizationHeaderPayload, payload)
		ctx.Next()

	}
}
