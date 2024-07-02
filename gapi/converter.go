package gapi

import (
	db "github.com/huzaifac137/simple_bank/db/sqlc"
	"github.com/huzaifac137/simple_bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user *db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
