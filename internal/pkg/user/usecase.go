package user

import (
	"mime/multipart"
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

type UseCase interface {
	Authorize(user *models.LoginUser) error
	SetAvatar(userId uint64, avatar multipart.File) error
	GetAvatar(userId uint64) (*os.File, error)
}
