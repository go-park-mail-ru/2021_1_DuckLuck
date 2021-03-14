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

func TestLocalRepository_AddProfile(t *testing.T) {
	rep := NewSessionLocalRepository()

	res, err = rep.AddProfile(userForAdd)
	require.NoError(t, err)
	require.Equal(t, res, retUser)
}

func TestLocalRepository_SelectProfileByEmail(t *testing.T) {
	rep := NewSessionLocalRepository()
	rep.AddProfile(userForAdd)

	res, err = rep.SelectProfileByEmail(badUserEmail)
	require.Error(t, server_errors.ErrUserNotFound)
	require.Nil(t, res)

	res, err = rep.SelectProfileByEmail(goodUserEmail)
	require.NoError(t, err)
	require.Equal(t, res, retUser)
}

func TestLocalRepository_SelectProfileById(t *testing.T) {
	rep := NewSessionLocalRepository()
	rep.AddProfile(userForAdd)

	res, err = rep.SelectProfileById(badUserId)
	require.Error(t, server_errors.ErrUserNotFound)
	require.Nil(t, res)

	res, err = rep.SelectProfileById(goodUserId)
	require.NoError(t, err)
	require.Equal(t, res, retUser)
}

func TestLocalRepository_UpdateProfile(t *testing.T) {
	rep := NewSessionLocalRepository()
	rep.AddProfile(userForAdd)

	err = rep.UpdateProfile(goodUserId, updateUser)
	require.NoError(t, err)
}

func TestLocalRepository_UpdateAvatar(t *testing.T) {
	rep := NewSessionLocalRepository()
	rep.AddProfile(userForAdd)

	err = rep.UpdateAvatar(badUserId, newAvatarName)
	require.Error(t, server_errors.ErrUserNotFound)

	err = rep.UpdateAvatar(goodUserId, newAvatarName)
	require.NoError(t, err)
}
