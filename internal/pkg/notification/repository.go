package notification

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type Repository interface {
	AddSubscribeUser(userId uint64, credentials *models.NotificationCredentials) error
	DeleteSubscribeUser(userId uint64) error
	SelectCredentialsByUserId(userId uint64) (*models.NotificationCredentials, error)
}
