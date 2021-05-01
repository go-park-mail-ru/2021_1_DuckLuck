package session

import "github.com/go-park-mail-ru/2021_1_DuckLuck/services/session/pkg/models"

type Repository interface {
	AddSession(session *models.Session) error
	SelectUserIdBySession(sessionValue string) (uint64, error)
	DeleteSessionByValue(sessionValue string) error
}
