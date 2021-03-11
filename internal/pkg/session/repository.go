package session

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session Repository

type Repository interface {
	Add(session *models.Session) error
	GetByValue(sessionValue string) (*models.Session, error)
	DestroyByValue(sessionValue string) error
}
