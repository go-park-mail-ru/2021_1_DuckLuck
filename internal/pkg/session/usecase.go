package session

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	Check(sessionCookieValue string) (*models.Session, error)
	Create(userId uint64) (*models.Session, error)
	DestroyCurrent(sessionCookieValue string) error
}
