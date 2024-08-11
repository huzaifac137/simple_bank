package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	SessionId             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	Username              string    `json:"username" binding:"required,alphanum"`
	Email                 string    `json:"email" binding:"required,email"`
	FullName              string    `json:"full_name" binding:"required"`
	CreatedAt             string    `json:"created_at"`
	PasswordChangedAt     string    `json:"password_changed_at"`
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

	maker, errr := token.NewJWTMaker(s.symmetricKey)
	if errr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessPayload, errrr := maker.CreateToken(req.Username, time.Minute*15)
	if errrr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := maker.CreateToken(req.Username, time.Hour*24)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.Id,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		fmt.Printf("error creating session: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := &LoginUserResponse{
		SessionId:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		Username:              user.Username,
		Email:                 user.Email,
		FullName:              user.FullName,
		CreatedAt:             user.CreatedAt.Format("2006-01-02 15:04:05"),
		PasswordChangedAt:     user.PasswordChangedAt.Format("2006-01-02 15:04:05"),
	}

	ctx.JSON(http.StatusOK, res)

}
