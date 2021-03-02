package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

type UseCase interface {
	Authorize(user *models.LoginUser) error
}
