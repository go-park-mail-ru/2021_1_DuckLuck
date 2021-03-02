package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type LocalRepository struct {
	data map[string]*models.Session
	mu   *sync.RWMutex
}

func NewSessionLocalRepository() session.Repository {
	return &LocalRepository{
		data: make(map[string]*models.Session, 7),
		mu:   &sync.RWMutex{},
	}
}

func (lr *LocalRepository) Add(session *models.Session) error {
	lr.mu.Lock()
	lr.data[session.Value] = session
	lr.mu.Unlock()

	return nil
}

func (lr *LocalRepository) GetByValue(sessionValue string) (*models.Session, error) {
	lr.mu.RLock()
	sess, ok := lr.data[sessionValue]
	lr.mu.Unlock()

	if !ok {
		return nil, server_errors.ErrSessionNotFound
	}

	return sess, nil
}

func (lr *LocalRepository) DestroyByValue(sessionValue string) error {
	lr.mu.Lock()
	delete(lr.data, sessionValue)
	lr.mu.Unlock()

	return nil
}
