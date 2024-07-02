package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/huzaifac137/simple_bank/db/sqlc"
	"github.com/huzaifac137/simple_bank/pb"
	"github.com/huzaifac137/simple_bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "Something went wrong: %s", err)
	}

	er := util.CheckPasswordHash(req.Password, user.HashedPassword)
	if er != nil {

		return nil, status.Errorf(codes.NotFound, "password is incorrect")
	}

	accessToken, accessPayload, errrr := s.tokenMaker.CreateToken(req.GetUsername(), time.Minute*15)
	if errrr != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong: %s", errrr)
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(req.GetUsername(), time.Hour*24)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong: %s", err)
	}

	metadata := s.ExtractMetadata(ctx)

	fmt.Printf("metadata: %+v\n", metadata)

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.Id,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong: %s", err)
	}

	res := &pb.LoginUserResponse{
		User:                  ConvertUser(&user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return res, nil
}
