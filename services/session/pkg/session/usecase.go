package session

import "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/models"

type UseCase interface {
	GetUserIdBySession(sessionCookieValue string) (uint64, error)
	CreateNewSession(userId uint64) (*models.Session, error)
	DestroySession(sessionCookieValue string) error
}
