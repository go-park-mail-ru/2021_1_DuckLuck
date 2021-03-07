package usecase

import (
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type UserUseCase struct {
	UserRepo user.Repository
}

func NewUseCase(userRepository user.Repository) user.UseCase {
	return &UserUseCase{
		UserRepo: userRepository,
	}
}

func (u *UserUseCase) Authorize(authUser *models.LoginUser) (*models.ProfileUser, error) {
	profileUser, err := u.UserRepo.GetByEmail(authUser.Email)
	if err != nil {
		return nil, server_errors.ErrIncorrectUserEmail
	}

	if profileUser.Password != authUser.Password {
		return nil, server_errors.ErrIncorrectUserPassword
	}

	return profileUser, nil
}

func (u *UserUseCase) SetAvatar(userId uint64, avatar string) (string, error) {
	// Destroy old user avatar
	profileUser, err := u.UserRepo.GetById(userId)
	if err == nil {
		err = os.Remove(configs.PathToUploadAvatar + profileUser.Avatar)
		return "", errors.ErrServerSystem
	}

	err = u.UserRepo.UpdateAvatar(userId, avatar)
	if err != nil {
		return "", err
	}

	return configs.UrlToAvatar + avatar, nil
}

func (u *UserUseCase) GetAvatar(userId uint64) (string, error) {
	profileUser, err := u.UserRepo.GetById(userId)
	if err != nil {
		return "", err
	}

	// If avatar not found -> return default_avatar.png
	var urlToFile string
	if _, err = os.Stat(configs.PathToUploadAvatar + profileUser.Avatar); err == nil {
		urlToFile = configs.UrlToAvatar + profileUser.Avatar

	} else {
		urlToFile = configs.UrlToAvatar + "default.png"
	}

	return urlToFile, nil
}
