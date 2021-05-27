package server

import (
	"context"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/grpc_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/pkg/user"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/proto/user"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/server/errors"
)

type AuthServer struct {
	UserUCase user.UseCase
}

func NewAuthServer(userUCase user.UseCase) *AuthServer {
	return &AuthServer{
		UserUCase: userUCase,
	}
}

func (s *AuthServer) Login(ctx context.Context,
	loginUser *proto.LoginUser) (*proto.UserId, error) {
	var err error
	defer func() {
		requireId := grpc_utils.MustGetRequireId(ctx)
		if err != nil {
			logger.LogError("auth_service_handler", "Login", requireId, err)
		}
	}()

	userId, err := s.UserUCase.Login(&models.LoginUser{
		Email:    loginUser.Email,
		Password: loginUser.Password,
	})

	if err != nil {
		return nil, errors.CreateError(err)
	}

	return &proto.UserId{Id: userId}, err
}

func (s *AuthServer) Signup(ctx context.Context,
	signupUser *proto.SignupUser) (*proto.UserId, error) {
	var err error
	defer func() {
		requireId := grpc_utils.MustGetRequireId(ctx)
		if err != nil {
			logger.LogError("auth_service_handler", "Signup", requireId, err)
		}
	}()

	userId, err := s.UserUCase.Signup(&models.SignupUser{
		Email:    signupUser.Email,
		Password: signupUser.Password,
	})

	if err != nil {
		return nil, errors.CreateError(err)
	}

	return &proto.UserId{Id: userId}, err
}
