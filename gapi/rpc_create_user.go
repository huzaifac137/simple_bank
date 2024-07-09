package gapi

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	"github.com/huzaifac137/simple_bank/grpc_util"
	"github.com/huzaifac137/simple_bank/pb"
	"github.com/huzaifac137/simple_bank/util"
	"github.com/huzaifac137/simple_bank/worker"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Validate the request
	voilations := createUserValidation(req)
	if voilations != nil {
		return nil, InvalidArgumentError(voilations)
	}

	hashedPass, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPass,
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "Username already exists: %s", err)

			default:
				return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
			}
		}

	}

	// send task to queue
	taskPayload := &worker.SendVerifyEmailPayload{Username: user.Username}
	options := []asynq.Option{}
	s.taskDistributer.DistributeSendVerifyEmail(ctx, taskPayload, options...)
	res := &pb.CreateUserResponse{
		User: ConvertUser(&user),
	}
	return res, nil
}

func createUserValidation(req *pb.CreateUserRequest) (voilations []*errdetails.BadRequest_FieldViolation) {
	err := grpc_util.ValidateUsername(req.GetUsername())
	if err != nil {
		voilations = append(voilations, FieldVoilation("username", err))
	}

	err = grpc_util.ValidatePassword(req.GetPassword())
	if err != nil {
		voilations = append(voilations, FieldVoilation("password", err))
	}

	err = grpc_util.ValidateEmail(req.GetEmail())
	if err != nil {
		voilations = append(voilations, FieldVoilation("email", err))
	}

	err = grpc_util.ValidateFullName(req.GetFullName())
	if err != nil {
		voilations = append(voilations, FieldVoilation("full_name", err))
	}

	return voilations

}
