package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user UseCase

type UseCase interface {
	Authorize(userRepo Repository, authUser *models.LoginUser) (*models.ProfileUser, error)
	UpdateProfile(userRepo Repository, userId uint64, updateUser *models.UpdateUser) error
	SetAvatar(userRepo Repository, userId uint64, avatar string) (string, error)
	GetAvatar(userRepo Repository, userId uint64) (string, error)
}
