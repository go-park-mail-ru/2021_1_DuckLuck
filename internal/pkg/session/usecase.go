package session

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session UseCase

type UseCase interface {
	GetUserIdBySession(sessionCookieValue string) (uint64, error)
	CreateNewSession(userId uint64) (*models.Session, error)
	DestroySession(sessionCookieValue string) error
}
