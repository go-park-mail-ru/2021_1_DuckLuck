package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user Repository

type Repository interface {
	AddProfile(user *models.ProfileUser) (uint64, error)
	UpdateProfile(userId uint64, user *models.UpdateUser) error
	SelectProfileById(userId uint64) (*models.ProfileUser, error)
	UpdateAvatar(userId uint64, avatarUrl string) error
}
