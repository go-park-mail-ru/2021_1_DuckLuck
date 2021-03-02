package usecase

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
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