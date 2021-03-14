package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type LocalRepository struct {
	data map[string]*models.Session
	mu   *sync.RWMutex
}

func NewSessionLocalRepository() session.Repository {
	return &LocalRepository{
		data: make(map[string]*models.Session, 0),
		mu:   &sync.RWMutex{},
	}
}

func (lr *LocalRepository) AddSession(session *models.Session) error {
	lr.mu.Lock()
	lr.data[session.Value] = session
	lr.mu.Unlock()

	return nil
}

func (lr *LocalRepository) SelectSessionByValue(sessionValue string) (*models.Session, error) {
	lr.mu.RLock()
	sess, ok := lr.data[sessionValue]
	lr.mu.RUnlock()

	if !ok {
		return nil, errors.ErrSessionNotFound
	}

	return sess, nil
}

func (lr *LocalRepository) DeleteByValue(sessionValue string) error {
	lr.mu.Lock()
	delete(lr.data, sessionValue)
	lr.mu.Unlock()

	return nil
}
