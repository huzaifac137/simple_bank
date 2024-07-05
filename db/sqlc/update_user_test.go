package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/huzaifac137/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserOnlyFullName(t *testing.T) {
	user := CreateRandomUser(t)
	newFullName := util.RandomString(10)
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.NotEqual(t, user.FullName, updatedUser.FullName)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)

}

func TestUpdateUserOnlyHashedPassowrd(t *testing.T) {
	user := CreateRandomUser(t)
	newHashedPassword, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, user.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.FullName, updatedUser.FullName)

}

func TestUpdateUserOnlyEmail(t *testing.T) {
	user := CreateRandomUser(t)
	newEmail := util.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, user.Email, updatedUser.Email)
	require.Equal(t, user.FullName, updatedUser.FullName)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)

}

func TestUpdateUserAllFields(t *testing.T) {
	user := CreateRandomUser(t)
	newEmail := util.RandomEmail()
	newFullName := util.RandomString(10)
	newHashedPassword, err := util.HashPassword(util.RandomString(10))
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, user.Email, updatedUser.Email)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.NotEqual(t, user.FullName, updatedUser.FullName)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, user.HashedPassword, updatedUser.HashedPassword)

}
