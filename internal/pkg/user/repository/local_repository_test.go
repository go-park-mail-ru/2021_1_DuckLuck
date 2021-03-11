package repository

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

const goodUserEmail = "goodEmail"
const badUserEmail = "badEmail"
const goodUserId = 0
const badUserId = 999
const newAvatarName = "newAvatar"

var userForAdd = &models.SignupUser{
	Email:    goodUserEmail,
	Password: "password",
}

var retUser = &models.ProfileUser{
	Id:       goodUserId,
	Email:    goodUserEmail,
	Password: "password",
}

var updateUser = &models.UpdateUser{
	FirstName: "firstName",
	LastName:  "lastName",
}

var res *models.ProfileUser
var err error

func TestLocalRepository_Add(t *testing.T) {
	rep := NewSessionLocalRepository()

	res, err = rep.Add(userForAdd)
	require.NoError(t, err)
	require.Equal(t, res, retUser)
}

func TestLocalRepository_GetByEmail(t *testing.T) {
	rep := NewSessionLocalRepository()
	rep.Add(userForAdd)

	res, err = rep.GetByEmail(badUserEmail)
	require.Error(t, server_errors.ErrUserNotFound)
	require.Nil(t, res)

	res, err = rep.GetByEmail(goodUserEmail)
	require.NoError(t, err)
	require.Equal(t, res, retUser)
}

func TestLocalRepository_GetById(t *testing.T) {
	rep := NewSessionLocalRepository()
	rep.Add(userForAdd)

	res, err = rep.GetById(badUserId)
	require.Error(t, server_errors.ErrUserNotFound)
	require.Nil(t, res)

	res, err = rep.GetById(goodUserId)
	require.NoError(t, err)
	require.Equal(t, res, retUser)
}

func TestLocalRepository_UpdateProfile(t *testing.T) {
	rep := NewSessionLocalRepository()
	rep.Add(userForAdd)

	err = rep.UpdateProfile(goodUserId, updateUser)
	require.NoError(t, err)
}

func TestLocalRepository_UpdateAvatar(t *testing.T) {
	rep := NewSessionLocalRepository()
	rep.Add(userForAdd)

	err = rep.UpdateAvatar(badUserId, newAvatarName)
	require.Error(t, server_errors.ErrUserNotFound)

	err = rep.UpdateAvatar(goodUserId, newAvatarName)
	require.NoError(t, err)
}
