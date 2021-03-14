package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user UseCase

type UseCase interface {
	Authorize(authUser *models.LoginUser) (*models.ProfileUser, error)
	UpdateProfile(userId uint64, updateUser *models.UpdateUser) error
	SetAvatar(userId uint64, avatar string) (string, error)
	GetAvatar(userId uint64) (string, error)
	GetUserById(userId uint64) (*models.ProfileUser, error)
	AddUser(user *models.SignupUser) (*models.ProfileUser, error)
}
