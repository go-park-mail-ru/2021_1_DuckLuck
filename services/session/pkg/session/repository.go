package session

import "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/session Repository

type Repository interface {
	AddSession(session *models.Session) error
	SelectUserIdBySession(sessionValue string) (uint64, error)
	DeleteSessionByValue(sessionValue string) error
}
