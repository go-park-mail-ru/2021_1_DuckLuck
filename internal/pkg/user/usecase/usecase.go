package usecase

import (
	"os"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/configs"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/password_hasher"
)

type UserUseCase struct {
	UserRepo user.Repository
}

func NewUseCase(repo user.Repository) user.UseCase {
	return &UserUseCase{
		UserRepo: repo,
	}
}

// Auth user
func (u *UserUseCase) Authorize(authUser *models.LoginUser) (*models.ProfileUser, error) {
	profileUser, err := u.UserRepo.SelectProfileByEmail(authUser.Email)
	if err != nil {
		return nil, errors.ErrIncorrectUserEmail
	}

	if ok := password_hasher.CompareHashAndPassword(profileUser.Password, authUser.Password); !ok {
		return nil, errors.ErrIncorrectUserPassword
	}

	return profileUser, nil
}

// Set new avatar
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

// Get user avatar
func (u *UserUseCase) GetAvatar(userId uint64) (string, error) {
	profileUser, err := u.UserRepo.SelectProfileById(userId)
	if err != nil {
		return "", err
	}

	return profileUser.Avatar.Url.String, nil
}

// Update user profile in repo
func (u *UserUseCase) UpdateProfile(userId uint64, updateUser *models.UpdateUser) error {
	return u.UserRepo.UpdateProfile(userId, updateUser)
}

// Get user profile by id
func (u *UserUseCase) GetUserById(userId uint64) (*models.ProfileUser, error) {
	return u.UserRepo.SelectProfileById(userId)
}

// Create new user in repo
func (u *UserUseCase) AddUser(user *models.SignupUser) (uint64, error) {
	if _, err := u.UserRepo.SelectProfileByEmail(user.Email); err == nil {
		return 0, errors.ErrEmailAlreadyExist
	}

	hashOfPassword, err := password_hasher.GenerateHashFromPassword(user.Password)
	if err != nil {
		return 0, errors.ErrHashFunctionFailed
	}

	return u.UserRepo.AddProfile(&models.ProfileUser{
		Email:    user.Email,
		Password: hashOfPassword,
	})
}
