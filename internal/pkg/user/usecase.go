package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

type UseCase interface {
	Authorize(userRepo Repository, authUser *models.LoginUser) (*models.ProfileUser, error)
	SetAvatar(userRepo Repository, userId uint64, avatar string) (string, error)
	GetAvatar(userRepo Repository, userId uint64) (string, error)
}
