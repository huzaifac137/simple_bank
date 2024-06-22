package token

import (
	"testing"
	"time"

	"github.com/huzaifac137/simple_bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	secretKey := util.RandomString(32)
	username := util.RandomOwner()
	issuedAt := time.Now()
	duration := time.Minute
	expiredAt := issuedAt.Add(duration)

	maker, err := NewPasetoMaker(secretKey)
	require.NoError(t, err)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.Id)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Minute)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Minute)

}
