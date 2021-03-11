package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

type Repository interface {
	Add(user *models.SignupUser) (*models.ProfileUser, error)
	UpdateProfile(userId uint64, user *models.UpdateUser) error
	GetByEmail(email string) (*models.ProfileUser, error)
	GetById(userId uint64) (*models.ProfileUser, error)
	UpdateAvatar(userId uint64, fileName string) error
}
