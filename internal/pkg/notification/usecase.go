package notification

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/notification UseCase

type UseCase interface {
	SubscribeUser(userId uint64, credentials *models.NotificationCredentials) error
	UnsubscribeUser(userId uint64, endpoint string) error
}
