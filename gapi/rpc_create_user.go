package gapi

import (
	"context"

	db "github.com/huzaifac137/simple_bank/db/sqlc"
	"github.com/huzaifac137/simple_bank/pb"
	"github.com/huzaifac137/simple_bank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
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
	res := &pb.CreateUserResponse{
		User: ConvertUser(&user),
	}
	return res, nil
}
