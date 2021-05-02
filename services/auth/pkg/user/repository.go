package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/models"
)

type Repository interface {
	AddProfile(user *models.AuthUser) (uint64, error)
	SelectUserByEmail(email string) (*models.AuthUser, error)
}
