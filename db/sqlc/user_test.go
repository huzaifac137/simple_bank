package db

import (
	"context"
	"testing"
	"time"

	"github.com/huzaifac137/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		t.Fatal(err)
	}
	require.NotEmpty(t, user)
	require.NoError(t, err)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}
func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	account1 := CreateRandomUser(t)
	account2, err := testQueries.GetUser(context.Background(), account1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.HashedPassword, account2.HashedPassword)
	require.Equal(t, account1.FullName, account2.FullName)
	require.Equal(t, account1.Email, account2.Email)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.WithinDuration(t, account1.PasswordChangedAt, account2.PasswordChangedAt, time.Second)

}
