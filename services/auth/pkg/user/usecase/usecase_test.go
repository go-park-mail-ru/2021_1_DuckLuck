package usecase

import (
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/user/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/server/tools/password_hasher"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUseCase_Login(t *testing.T) {
	loginUser := models.LoginUser{
		Email:    "test@test.ru",
		Password: "qwerty",
	}

	incorrectLoginUser := models.LoginUser{
		Email:    "test@test.ru",
		Password: "qwer",
	}

	hashedPassword, _ := password_hasher.GenerateHashFromPassword("qwerty")
	profileUser := models.AuthUser{
		Id:       1,
		Password: hashedPassword,
		Email:    "test@test.ru",
	}

	userId := uint64(1)

	t.Run("Login_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectUserByEmail(loginUser.Email).
			Return(&profileUser, nil)

		userUCase := NewUseCase(userRep)

		returnedId, err := userUCase.Login(&loginUser)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, returnedId, userId, "not equal user ids")
	})

	t.Run("Login_incorrect_password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectUserByEmail(loginUser.Email).
			Return(&profileUser, nil)

		userUCase := NewUseCase(userRep)

		_, err := userUCase.Login(&incorrectLoginUser)
		assert.Equal(t, err, errors.ErrIncorrectAuthData, "expected error")
	})

	t.Run("Login_incorrect_email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectUserByEmail(loginUser.Email).
			Return(nil, errors.ErrDBInternalError)

		userUCase := NewUseCase(userRep)

		_, err := userUCase.Login(&incorrectLoginUser)
		assert.Equal(t, err, errors.ErrIncorrectAuthData, "expected error")
	})
}
