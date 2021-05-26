package usecase

import (
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	notification_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/notification/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNotificationUseCase_SubscribeUser(t *testing.T) {
	credentials := models.NotificationCredentials{}
	userId := uint64(3)
	subscribes := models.Subscribes{Credentials: make(map[string]*models.NotificationKeys)}

	t.Run("SubscribeUser_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationRepo := notification_mock.NewMockRepository(ctrl)
		notificationRepo.
			EXPECT().
			SelectCredentialsByUserId(userId).
			Return(&subscribes, nil)

		notificationRepo.
			EXPECT().
			AddSubscribeUser(userId, &subscribes).
			Return(nil)

		notificationUCase := NewUseCase(notificationRepo)

		err := notificationUCase.SubscribeUser(userId, &credentials)
		assert.NoError(t, err, "unexpected error")
	})
}

func TestNotificationUseCase_UnsubscribeUser(t *testing.T) {
	credentials := models.NotificationCredentials{}
	userId := uint64(3)
	subscribes := models.Subscribes{
		Credentials: map[string]*models.NotificationKeys{
			"test": {},
		},
	}
	manySubscribes := models.Subscribes{
		Credentials: map[string]*models.NotificationKeys{
			"test":  {},
			"test2": {},
		},
	}

	t.Run("UnsubscribeUser_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationRepo := notification_mock.NewMockRepository(ctrl)
		notificationRepo.
			EXPECT().
			SelectCredentialsByUserId(userId).
			Return(&subscribes, nil)

		notificationRepo.
			EXPECT().
			DeleteSubscribeUser(userId).
			Return(nil)

		notificationUCase := NewUseCase(notificationRepo)

		err := notificationUCase.UnsubscribeUser(userId, credentials.Endpoint)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("UnsubscribeUser_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationRepo := notification_mock.NewMockRepository(ctrl)
		notificationRepo.
			EXPECT().
			SelectCredentialsByUserId(userId).
			Return(&manySubscribes, nil)

		notificationRepo.
			EXPECT().
			AddSubscribeUser(userId, &manySubscribes).
			Return(nil)

		notificationUCase := NewUseCase(notificationRepo)

		err := notificationUCase.UnsubscribeUser(userId, credentials.Endpoint)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("UnsubscribeUser_notification_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationRepo := notification_mock.NewMockRepository(ctrl)
		notificationRepo.
			EXPECT().
			SelectCredentialsByUserId(userId).
			Return(nil, errors.ErrInternalError)

		notificationUCase := NewUseCase(notificationRepo)

		err := notificationUCase.UnsubscribeUser(userId, credentials.Endpoint)
		assert.Error(t, err, "unexpected error")
	})
}
