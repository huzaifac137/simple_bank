package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	token "github.com/huzaifac137/simple_bank/token"
	"github.com/huzaifac137/simple_bank/util"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

type CreateUserResponse struct {
	Username          string `json:"username" binding:"required,alphanum"`
	Password          string `json:"password" binding:"required,min=6"`
	Email             string `json:"email" binding:"required,email"`
	FullName          string `json:"full_name" binding:"required"`
	CreatedAt         string `json:"created_at"`
	PasswordChangedAt string `json:"password_changed_at"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPass, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPass,
		Email:          req.Email,
		FullName:       req.FullName,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			default:
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := &CreateUserResponse{
		Username:          user.Username,
		Password:          user.HashedPassword,
		Email:             user.Email,
		FullName:          user.FullName,
		CreatedAt:         user.CreatedAt.Format("2006-01-02 15:04:05"),
		PasswordChangedAt: user.PasswordChangedAt.Format("2006-01-02 15:04:05"),
	}
	ctx.JSON(http.StatusOK, response)
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken       string `json:"access_token"`
	Username          string `json:"username" binding:"required,alphanum"`
	Email             string `json:"email" binding:"required,email"`
	FullName          string `json:"full_name" binding:"required"`
	CreatedAt         string `json:"created_at"`
	PasswordChangedAt string `json:"password_changed_at"`
}

func (s *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	er := util.CheckPasswordHash(req.Password, user.HashedPassword)
	if er != nil {
		passwordErr := errors.New("password is incorrect")
		ctx.JSON(http.StatusBadRequest, errorResponse(passwordErr))
		return
	}

	maker, errr := token.NewPasetoMaker(s.symmetricKey)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, errrr := maker.CreateToken(req.Username, time.Minute)
	if errrr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := &LoginUserResponse{
		AccessToken:       token,
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		CreatedAt:         user.CreatedAt.Format("2006-01-02 15:04:05"),
		PasswordChangedAt: user.PasswordChangedAt.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, res)

}
