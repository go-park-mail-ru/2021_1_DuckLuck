package usecase

import (
	"mime/multipart"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/hasher"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/s3_utils"
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
		return nil, errors.ErrIncorrectAuthData
	}

	if ok := hasher.CompareHashAndPassword(profileUser.Password, authUser.Password); !ok {
		return nil, errors.ErrIncorrectAuthData
	}

	return profileUser, nil
}

// Set new avatar
func (u *UserUseCase) SetAvatar(userId uint64, file *multipart.File, header *multipart.FileHeader) (string, error) {
	// Upload new user avatar to S3
	fileName, err := s3_utils.UploadMultipartFile("avatar", file, header)
	if err != nil {
		return "", err
	}

	// Destroy old user avatar
	profileUser, err := u.UserRepo.SelectProfileById(userId)
	if err == nil && profileUser.Avatar.Url != "" {
		if err = s3_utils.DeleteFile(profileUser.Avatar.Url); err != nil {
			return "", err
		}
	}

	err = u.UserRepo.UpdateAvatar(userId, fileName)
	if err != nil {
		return "", err
	}

	return s3_utils.PathToFile(fileName), nil
}

// Get user avatar
func (u *UserUseCase) GetAvatar(userId uint64) (string, error) {
	profileUser, err := u.UserRepo.SelectProfileById(userId)
	if err != nil {
		return "", errors.ErrUserNotFound
	}

	return s3_utils.PathToFile(profileUser.Avatar.Url), nil
}

// Update user profile in repo
func (u *UserUseCase) UpdateProfile(userId uint64, updateUser *models.UpdateUser) error {
	return u.UserRepo.UpdateProfile(userId, updateUser)
}

// Get user profile by id
func (u *UserUseCase) GetUserById(userId uint64) (*models.ProfileUser, error) {
	userById, err := u.UserRepo.SelectProfileById(userId)
	if err != nil {
		return nil, err
	}
	userById.Avatar.Url = s3_utils.PathToFile(userById.Avatar.Url)
	return userById, nil
}

// Create new user in repo
func (u *UserUseCase) AddUser(user *models.SignupUser) (uint64, error) {
	if _, err := u.UserRepo.SelectProfileByEmail(user.Email); err == nil {
		return 0, errors.ErrEmailAlreadyExist
	}

	hashOfPassword, err := hasher.GenerateHashFromPassword(user.Password)
	if err != nil {
		return 0, errors.ErrHashFunctionFailed
	}

	return u.UserRepo.AddProfile(&models.ProfileUser{
		Email:    user.Email,
		Password: hashOfPassword,
	})
}
