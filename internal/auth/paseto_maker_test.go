package auth

import (
	"testing"
	"time"

	"github.com/SecureParadise/go_attendence/internal/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	symmetricKey := util.RandomString(32)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	username := util.RandomOwner()
	role := "student"
	duration := time.Minute
	tokenType := TokenType('A')

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration, tokenType)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, tokenType)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.Equal(t, tokenType, payload.Type)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	symmetricKey := util.RandomString(32)
	maker, err := NewPasetoMaker(symmetricKey)
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), "student", -time.Minute, TokenType('A'))
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, TokenType('A'))
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
