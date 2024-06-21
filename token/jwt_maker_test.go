package token

import (
	"testing"
	"time"

	"github.com/huzaifac137/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	secretKey := util.RandomString(32)
	username := util.RandomOwner()
	issuedAt := time.Now()
	duration := time.Minute
	expiredAt := issuedAt.Add(duration)

	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	token, err := maker.createToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.verifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Minute)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Minute)

}

func TestJWTMakerExpired(t *testing.T) {
	secretKey := util.RandomString(32)
	username := util.RandomOwner()
	duration := -time.Minute

	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	token, err := maker.createToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.verifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())

}
