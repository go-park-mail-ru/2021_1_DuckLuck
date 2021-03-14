package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/mock"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

const goodId = uint64(1)
const badId = uint64(999)
const goodEmail = "goodEmail"
const goodPassword = "goodPassword"
const badEmail = "badEmail"
const badPassword = "badPassword"
const goodAvatar = "avatar"

var retUser = &models.ProfileUser{
	Id:       goodId,
	Password: goodPassword,
	Email:    goodEmail,
	Avatar:   models.Avatar{
		Url: goodAvatar,
	},
}
var err error

func TestUserUseCase_Authorize(t *testing.T) {
	authUser := &models.LoginUser{
		Email:    goodEmail,
		Password: goodPassword,
	}

	useCase := NewUseCase()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockRepository.EXPECT().GetByEmail(goodEmail).Times(2).Return(retUser, nil)
	mockRepository.EXPECT().GetByEmail(badEmail).Times(1).Return(nil, server_errors.ErrIncorrectUserEmail)

	var res *models.ProfileUser

	res, err = useCase.Authorize(mockRepository, authUser)
	require.NoError(t, err)
	require.Equal(t, res, retUser)

	authUser.Password = badPassword
	res, err = useCase.Authorize(mockRepository, authUser)
	require.Error(t, err, server_errors.ErrIncorrectUserPassword)
	require.Nil(t, res)
	authUser.Password = goodPassword

	authUser.Email = badEmail
	res, err = useCase.Authorize(mockRepository, authUser)
	require.Error(t, err, server_errors.ErrIncorrectUserEmail)
	require.Nil(t, res)
}

func TestUserUseCase_SetAvatar(t *testing.T) {
	useCase := NewUseCase()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockRepository.EXPECT().GetById(goodId).Times(1).Return(retUser, nil)
	mockRepository.EXPECT().UpdateAvatar(goodId, configs.UrlToAvatar+goodAvatar).Times(1).Return(nil)

	var res string

	res, err = useCase.SetAvatar(mockRepository, goodId, goodAvatar)
	require.NoError(t, err)
	require.Equal(t, res, configs.UrlToAvatar+goodAvatar)
}

func TestUserUseCase_GetAvatar(t *testing.T) {
	useCase := NewUseCase()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockRepository.EXPECT().GetById(goodId).Times(1).Return(retUser, nil)
	mockRepository.EXPECT().GetById(badId).Times(1).Return(nil, server_errors.ErrUserNotFound)

	var res string

	res, err = useCase.GetAvatar(mockRepository, goodId)
	require.NoError(t, err)
	require.Equal(t, res, goodAvatar)

	res, err = useCase.GetAvatar(mockRepository, badId)
	require.Error(t, err, server_errors.ErrUserNotFound)
	require.Equal(t, res, "")
}

func TestUserUseCase_UpdateProfile(t *testing.T) {
	updateUser := &models.UpdateUser{
		FirstName: "newFirstName",
		LastName:  "newLastName",
	}

	useCase := NewUseCase()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock.NewMockRepository(ctrl)
	mockRepository.EXPECT().UpdateProfile(goodId, updateUser).Times(1).Return(nil)
	mockRepository.EXPECT().UpdateProfile(badId, updateUser).Times(1).Return(server_errors.ErrUserNotFound)

	err = useCase.UpdateProfile(mockRepository, goodId, updateUser)
	require.NoError(t, err)

	err = useCase.UpdateProfile(mockRepository, badId, updateUser)
	require.Error(t, err, server_errors.ErrUserNotFound)
}
