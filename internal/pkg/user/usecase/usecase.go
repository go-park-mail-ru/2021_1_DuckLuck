package usecase

import (
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	UserRepo user.Repository
}

func NewUseCase(repo user.Repository) user.UseCase {
	return &UserUseCase{
		UserRepo: repo,
	}
}

func (u *UserUseCase) Authorize(authUser *models.LoginUser) (*models.ProfileUser, error) {
	profileUser, err := u.UserRepo.SelectProfileByEmail(authUser.Email)
	if err != nil {
		return nil, errors.ErrIncorrectUserEmail
	}

	if err = bcrypt.CompareHashAndPassword([]byte(profileUser.Password), []byte(authUser.Password)); err != nil {
		return nil, errors.ErrIncorrectUserPassword
	}

	return profileUser, nil
}

func (u *UserUseCase) SetAvatar(userId uint64, avatar string) (string, error) {
	// Destroy old user avatar
	profileUser, err := u.UserRepo.SelectProfileById(userId)
	if err == nil || profileUser.Avatar.Url.Valid {
		err = os.Remove(configs.PathToUpload + profileUser.Avatar.Url.String)
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

	return profileUser.Avatar.Url.String, nil
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

func (u *UserUseCase) AddUser(user *models.SignupUser) (uint64, error) {
	if _, err := u.UserRepo.SelectProfileByEmail(user.Email); err == nil {
		return 0, errors.ErrEmailAlreadyExist
	}

	hashOfPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return 0, errors.ErrHashFunctionFailed
	}

	return u.UserRepo.AddProfile(&models.SignupUser{
		Email:    user.Email,
		Password: string(hashOfPassword),
	})
}
