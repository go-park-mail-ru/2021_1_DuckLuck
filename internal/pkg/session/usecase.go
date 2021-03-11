package session

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session UseCase

type UseCase interface {
	Check(sessionCookieValue string) (*models.Session, error)
	Create(userId uint64) (*models.Session, error)
	DestroyCurrent(sessionCookieValue string) error
}
