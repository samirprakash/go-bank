package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedhPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedhPassword1)

	err = CheckPassword(password, hashedhPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedhPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedhPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedhPassword2)
	require.NotEqual(t, hashedhPassword1, hashedhPassword2)
}
