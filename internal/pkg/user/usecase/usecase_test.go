package usecase

import (
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/hasher"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/s3_utils"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUseCase_Authorize(t *testing.T) {
	loginUser := models.LoginUser{
		Email:    "test@test.ru",
		Password: "qwerty",
	}

	incorrectLoginUser := models.LoginUser{
		Email:    "test@test.ru",
		Password: "qwer",
	}

	hashedPassword, _ := hasher.GenerateHashFromPassword("qwerty")
	profileUser := models.ProfileUser{
		Id:        1,
		FirstName: "test",
		LastName:  "last",
		Email:     "test@test.ru",
		Password:  hashedPassword,
		Avatar: models.Avatar{
			Url: "httt://test.png",
		},
	}

	t.Run("Authorize_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectProfileByEmail(loginUser.Email).
			Return(&profileUser, nil)

		userUCase := NewUseCase(userRep)

		userData, err := userUCase.Authorize(&loginUser)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, *userData, profileUser, "not equal data")
	})

	t.Run("Authorize_incorrect_password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectProfileByEmail(loginUser.Email).
			Return(&profileUser, nil)

		userUCase := NewUseCase(userRep)

		_, err := userUCase.Authorize(&incorrectLoginUser)
		assert.Equal(t, err, errors.ErrIncorrectAuthData, "expected error")
	})

	t.Run("Authorize_incorrect_email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectProfileByEmail(loginUser.Email).
			Return(nil, errors.ErrDBInternalError)

		userUCase := NewUseCase(userRep)

		_, err := userUCase.Authorize(&incorrectLoginUser)
		assert.Equal(t, err, errors.ErrIncorrectAuthData, "expected error")
	})
}

func TestUserUseCase_GetAvatar(t *testing.T) {
	profileUser := models.ProfileUser{
		Id:        1,
		FirstName: "test",
		LastName:  "last",
		Email:     "test@test.ru",
		Password:  []byte{1, 2, 3, 43, 32},
		Avatar: models.Avatar{
			Url: "test.png",
		},
	}
	fileUrl := s3_utils.PathToFile(profileUser.Avatar.Url)

	t.Run("GetAvatar_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectProfileById(profileUser.Id).
			Return(&profileUser, nil)

		userUCase := NewUseCase(userRep)

		url, err := userUCase.GetAvatar(profileUser.Id)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, fileUrl, url, "not equal data")
	})

	t.Run("GetAvatar_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectProfileById(profileUser.Id).
			Return(nil, errors.ErrDBInternalError)

		userUCase := NewUseCase(userRep)

		_, err := userUCase.GetAvatar(profileUser.Id)
		assert.Equal(t, errors.ErrUserNotFound, err, "not equal errors")
	})
}

func TestUserUseCase_UpdateProfile(t *testing.T) {
	updateUser := models.UpdateUser{
		FirstName: "test",
		LastName:  "last test",
	}
	userId := uint64(3)

	t.Run("UpdateProfile_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			UpdateProfile(userId, &updateUser).
			Return(nil)

		userUCase := NewUseCase(userRep)

		err := userUCase.UpdateProfile(userId, &updateUser)
		assert.NoError(t, err, "unexpected error")
	})
}

func TestUserUseCase_GetUserById(t *testing.T) {
	profileUser := models.ProfileUser{
		Id:        3,
		FirstName: "test",
		LastName:  "last",
		Email:     "test@test.ru",
		Password:  []byte{1, 2, 3, 43, 32},
		Avatar: models.Avatar{
			Url: "test.png",
		},
	}
	userId := uint64(3)

	t.Run("UpdateProfile_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectProfileById(userId).
			Return(&profileUser, nil)

		userUCase := NewUseCase(userRep)

		data, err := userUCase.GetUserById(userId)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, *data, profileUser, "not equal data")
	})
}

func TestUserUseCase_AddUser(t *testing.T) {
	signupUser := models.SignupUser{
		Email:    "test@test.ru",
		Password: "qwerty",
	}

	hashedPassword, _ := hasher.GenerateHashFromPassword("qwerty")
	profileUser := models.ProfileUser{
		Id:        1,
		FirstName: "test",
		LastName:  "last",
		Email:     "test@test.ru",
		Password:  hashedPassword,
		Avatar: models.Avatar{
			Url: "httt://test.png",
		},
	}

	t.Run("AddUser_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectProfileByEmail(profileUser.Email).
			Return(&profileUser, nil)

		userUCase := NewUseCase(userRep)

		_, err := userUCase.AddUser(&signupUser)
		assert.Equal(t, errors.ErrEmailAlreadyExist, err, "not equal errors")
	})
}
