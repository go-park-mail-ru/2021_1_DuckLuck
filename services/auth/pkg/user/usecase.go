package user

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/models"
)

type UseCase interface {
	Login(loginUser *models.LoginUser) (uint64, error)
	Signup(signupUser *models.SignupUser) (uint64, error)
}
