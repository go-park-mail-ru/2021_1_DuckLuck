package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/models"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/user UseCase

type UseCase interface {
	Login(loginUser *models.LoginUser) (uint64, error)
	Signup(signupUser *models.SignupUser) (uint64, error)
}
