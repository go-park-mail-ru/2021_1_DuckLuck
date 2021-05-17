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
	return u.NotificationRepo.AddSubscribeUser(userId, credentials)
}

func (u *NotificationUseCase) UnsubscribeUser(userId uint64) error {
	return u.NotificationRepo.DeleteSubscribeUser(userId)
}
