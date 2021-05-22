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
	var subscribes *models.Subscribes
	subscribes, err := u.NotificationRepo.SelectCredentialsByUserId(userId)
	if err != nil || subscribes.Credentials == nil || subscribes == nil {
		subscribes = &models.Subscribes{}
		subscribes.Credentials = make(map[string]*models.NotificationKeys, 0)
	}

	subscribes.Credentials[credentials.Endpoint] = &credentials.Keys
	return u.NotificationRepo.AddSubscribeUser(userId, subscribes)
}

func (u *NotificationUseCase) UnsubscribeUser(userId uint64, endpoint string) error {
	userSubscribes, err := u.NotificationRepo.SelectCredentialsByUserId(userId)
	if err != nil {
		return err
	}

	if len(userSubscribes.Credentials) == 1 {
		return u.NotificationRepo.DeleteSubscribeUser(userId)
	}

	delete(userSubscribes.Credentials, endpoint)
	return u.NotificationRepo.AddSubscribeUser(userId, userSubscribes)
}
