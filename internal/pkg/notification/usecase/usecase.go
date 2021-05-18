package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/notification"
)

type NotificationUseCase struct {
	NotificationRepo notification.Repository
}

func NewUseCase(notificationRepo notification.Repository) notification.UseCase {
	return &NotificationUseCase{
		NotificationRepo: notificationRepo,
	}
}

func (u *NotificationUseCase) SubscribeUser(userId uint64,
	credentials *models.NotificationCredentials) error {
	userSubscribes, err := u.NotificationRepo.SelectCredentialsByUserId(userId)
	if err != nil || userSubscribes.Credentials == nil || userSubscribes == nil {
		subscribes := &models.Subscribes{}
		subscribes.Credentials = make(map[string]*models.NotificationKeys, 0)
		return u.NotificationRepo.AddSubscribeUser(userId, subscribes)
	}
	userSubscribes.Credentials[credentials.Endpoint] = &credentials.Keys

	return u.NotificationRepo.AddSubscribeUser(userId, userSubscribes)
}

func (u *NotificationUseCase) UnsubscribeUser(userId uint64, endpoint string) error {
	userSubscribes, err := u.NotificationRepo.SelectCredentialsByUserId(userId)
	if err != nil {
		return err
	}

	if len(userSubscribes.Credentials) == 1 {
		if err = u.NotificationRepo.DeleteSubscribeUser(userId); err != nil {
			return err
		}
	}

	delete(userSubscribes.Credentials, endpoint)
	return u.NotificationRepo.AddSubscribeUser(userId, userSubscribes)
}
