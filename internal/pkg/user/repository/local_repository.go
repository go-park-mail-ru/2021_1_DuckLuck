package repository

import (
	"sync"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type LocalRepository struct {
	data map[uint64]*models.ProfileUser
	mu   *sync.RWMutex
}

func NewSessionLocalRepository() user.Repository {
	return &LocalRepository{
		data: make(map[uint64]*models.ProfileUser, 0),
		mu:   &sync.RWMutex{},
	}
}

func (lr *LocalRepository) AddProfile(user *models.SignupUser) (*models.ProfileUser, error) {
	newUser := &models.ProfileUser{
		Id:        uint64(len(lr.data)),
		FirstName: "",
		LastName:  "",
		Email:     user.Email,
		Password:  user.Password,
		Avatar: models.Avatar{
			Url: "",
		},
	}
	lr.mu.Lock()
	lr.data[newUser.Id] = newUser
	lr.mu.Unlock()

	return newUser, nil
}

func (lr *LocalRepository) SelectProfileByEmail(email string) (*models.ProfileUser, error) {
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
		return nil, errors.ErrUserNotFound
	}

	return userByEmail, nil
}

func (lr *LocalRepository) SelectProfileById(userId uint64) (*models.ProfileUser, error) {
	lr.mu.RLock()
	userById, ok := lr.data[userId]
	lr.mu.RUnlock()

	if !ok {
		return nil, errors.ErrUserNotFound
	}

	return userById, nil
}

func (lr *LocalRepository) UpdateProfile(userId uint64, user *models.UpdateUser) error {
	lr.mu.Lock()
	lr.data[userId].FirstName = user.FirstName
	lr.data[userId].LastName = user.LastName
	lr.mu.Unlock()

	return nil
}

func (lr *LocalRepository) UpdateAvatar(userId uint64, fileName string) error {
	lr.mu.RLock()
	_, ok := lr.data[userId]
	lr.mu.RUnlock()

	if !ok {
		return errors.ErrUserNotFound
	}

	lr.mu.Lock()
	lr.data[userId].Avatar.Url = fileName
	lr.mu.Unlock()

	return nil
}
