package session

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

type Repository interface {
	Add(session *models.Session) error
	GetByValue(sessionValue string) (*models.Session, error)
	DestroyByValue(sessionValue string) error
}
