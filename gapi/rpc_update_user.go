package gapi

import (
	"context"
	"database/sql"

	db "github.com/huzaifac137/simple_bank/db/sqlc"
	"github.com/huzaifac137/simple_bank/grpc_util"
	"github.com/huzaifac137/simple_bank/pb"
	"github.com/huzaifac137/simple_bank/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// Validate the request
	voilations := UpdateUserValidation(req)
	if voilations != nil {
		return nil, InvalidArgumentError(voilations)
	}

	_, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "Something went wrong: %s", err)
	}

	var updatedUser db.UpdateUserParams
	updatedUser.Username = req.GetUsername()

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Something went wrong: %s", err)
		}
		updatedUser.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

	}

	if req.Email != nil {
		updatedUser.Email = sql.NullString{
			String: req.GetEmail(),
			Valid:  true,
		}
	}

	if req.FullName != nil {
		updatedUser.FullName = sql.NullString{
			String: req.GetFullName(),
			Valid:  true,
		}
	}

	resUser, err := s.store.UpdateUser(ctx, updatedUser)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Something went wrong: %s", err)
	}

	res := &pb.UpdateUserResponse{
		User: ConvertUser(&resUser),
	}

	return res, nil
}

func UpdateUserValidation(req *pb.UpdateUserRequest) (voilations []*errdetails.BadRequest_FieldViolation) {
	err := grpc_util.ValidateUsername(req.GetUsername())
	if err != nil {
		voilations = append(voilations, FieldVoilation("username", err))
	}

	if req.Password != nil {
		err = grpc_util.ValidatePassword(req.GetPassword())
		if err != nil {
			voilations = append(voilations, FieldVoilation("password", err))
		}
	}

	if req.Email != nil {
		err = grpc_util.ValidateEmail(req.GetEmail())
		if err != nil {
			voilations = append(voilations, FieldVoilation("email", err))
		}
	}

	if req.FullName != nil {
		err = grpc_util.ValidateFullName(req.GetFullName())
		if err != nil {
			voilations = append(voilations, FieldVoilation("full_name", err))
		}
	}

	return voilations

}
