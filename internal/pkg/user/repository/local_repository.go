package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type LocalRepository struct {
	data map[uint64]*models.ProfileUser
	mu   *sync.RWMutex
}

func NewSessionLocalRepository() user.Repository {
	return &LocalRepository{
		data: make(map[uint64]*models.ProfileUser, 7),
		mu:   &sync.RWMutex{},
	}
}

func (lr *LocalRepository) Add(user *models.SignupUser) (*models.ProfileUser, error) {
	newUser := &models.ProfileUser{
		Id:        uint64(len(lr.data)),
		FirstName: "",
		LastName:  "",
		Email:     user.Email,
		Password:  user.Password,
		Avatar:    "",
	}
	lr.mu.Lock()
	lr.data[newUser.Id] = newUser
	lr.mu.Unlock()

	return newUser, nil
}

func (lr *LocalRepository) GetByEmail(email string) (*models.ProfileUser, error) {
	lr.mu.RLock()
	ok := false
	var userByEmail *models.ProfileUser
	for _, currentUser := range lr.data {
		if currentUser.Email == email {
			ok = true
			userByEmail = currentUser
		}
	}
	lr.mu.RUnlock()

	if !ok {
		return nil, server_errors.ErrUserNotFound
	}

	return userByEmail, nil
}

func (lr *LocalRepository) GetById(userId uint64) (*models.ProfileUser, error) {
	lr.mu.RLock()
	userById, ok := lr.data[userId]
	lr.mu.Unlock()

	if !ok {
		return nil, server_errors.ErrUserNotFound
	}

	return userById, nil
}

func (lr *LocalRepository) Update(user *models.ProfileUser) error {
	lr.mu.RLock()
	lr.data[user.Id] = user
	lr.mu.Unlock()

	return nil
}
