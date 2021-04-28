package usecase

import (
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUseCase_GetUserIdBySession(t *testing.T) {
	sessionCookieValue := "test"
	userId := uint64(4)

	t.Run("GetUserIdBySession_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			SelectUserIdBySession(sessionCookieValue).
			Return(userId, nil)

		userUCase := NewUseCase(userRep)

		userData, err := userUCase.GetUserIdBySession(sessionCookieValue)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, userId, userData, "not equal data")
	})
}

func TestUserUseCase_DestroySession(t *testing.T) {
	sessionCookieValue := "test"

	t.Run("GetUserIdBySession_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRep := mock.NewMockRepository(ctrl)

		userRep.
			EXPECT().
			DeleteSessionByValue(sessionCookieValue).
			Return(nil)

		userUCase := NewUseCase(userRep)

		err := userUCase.DestroySession(sessionCookieValue)
		assert.NoError(t, err, "unexpected error")
	})
}
