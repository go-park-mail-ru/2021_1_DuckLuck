package usecase

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	uuid "github.com/satori/go.uuid"
)

type UserUseCase struct {
	UserRepo user.Repository
}

func NewUseCase(userRepository user.Repository) *UserUseCase {
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

func (u *UserUseCase) SetAvatar(userId uint64, avatar multipart.File) (string, error) {
	// Destroy old user avatar
	profileUser, err := u.UserRepo.GetById(userId)
	if err == nil {
		err = os.Remove(fmt.Sprintf("%s/%s.png", configs.PathToUploads, profileUser.Avatar))
		return "", errors.ErrServerSystem
	}

	newName := uuid.NewV4().String()
	f, err := os.Create(fmt.Sprintf("%s/%s.png", configs.PathToUploads, newName))
	if err != nil {
		return "", errors.ErrServerSystem
	}
	defer f.Close()

	io.Copy(f, avatar)

	err = u.UserRepo.UpdateAvatar(userId, newName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.png", newName), nil
}

func (u *UserUseCase) GetAvatar(userId uint64) (*os.File, error) {
	profileUser, err := u.UserRepo.GetById(userId)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fmt.Sprintf("%s/%s.png", configs.PathToUploads, profileUser.Avatar))
	// If avatar not found -> return default_avatar.png
	if err != nil {
		f, _ = os.Open(fmt.Sprintf("%s/default_avatar.png", configs.PathToUploads))
	}

	return f, nil
}
