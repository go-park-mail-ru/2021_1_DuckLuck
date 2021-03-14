package usecase

import (
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type UserUseCase struct{
	UserRepo  user.Repository
}

func NewUseCase(repo  user.Repository) user.UseCase {
	return &UserUseCase{
		UserRepo: repo,
	}
}

func (u *UserUseCase) Authorize(authUser *models.LoginUser) (*models.ProfileUser, error) {
	profileUser, err := u.UserRepo.SelectProfileByEmail(authUser.Email)
	if err != nil {
		return nil, errors.ErrIncorrectUserEmail
	}

	if profileUser.Password != authUser.Password {
		return nil, errors.ErrIncorrectUserPassword
	}

	return profileUser, nil
}

func (u *UserUseCase) SetAvatar(userId uint64, avatar string) (string, error) {
	// Destroy old user avatar
	profileUser, err := u.UserRepo.SelectProfileById(userId)
	if err == nil || profileUser.Avatar.Url != "" {
		err = os.Remove(configs.PathToUpload + profileUser.Avatar.Url)
	}

	err = u.UserRepo.UpdateAvatar(userId, configs.UrlToAvatar+avatar)
	if err != nil {
		return "", err
	}

	return configs.UrlToAvatar + avatar, nil
}

func (u *UserUseCase) GetAvatar(userId uint64) (string, error) {
	profileUser, err := u.UserRepo.SelectProfileById(userId)
	if err != nil {
		return "", err
	}

	return profileUser.Avatar.Url, nil
}

func (u *UserUseCase) UpdateProfile(userId uint64, updateUser *models.UpdateUser) error {
	err := u.UserRepo.UpdateProfile(userId, updateUser)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) GetUserById(userId uint64) (*models.ProfileUser, error) {
	return u.UserRepo.SelectProfileById(userId)
}

func (u *UserUseCase) AddUser(user *models.SignupUser) (*models.ProfileUser, error) {
	return u.UserRepo.AddProfile(user)
}
