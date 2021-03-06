package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

type UseCase interface {
	Authorize(user *models.LoginUser) error
	SetAvatar(userId uint64, avatar string) (string, error)
	GetAvatar(userId uint64) (string, error)
}
