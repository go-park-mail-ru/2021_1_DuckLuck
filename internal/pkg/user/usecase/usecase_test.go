package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	user_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	auth_service "github.com/go-park-mail-ru/2021_1_DuckLuck/services/auth/proto/user"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUseCase_Authorize(t *testing.T) {
	reqLoginUser := &auth_service.LoginUser{
		Email:    "test",
		Password: "name",
	}
	loginUser := &models.LoginUser{
		Email:    "test",
		Password: "name",
	}

	userId := &auth_service.UserId{Id: 12}

	t.Run("Authorize_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		authClient := auth_service.NewMockAuthServiceClient(ctrl)

		authClient.
			EXPECT().
			Login(context.Background(), reqLoginUser).
			Return(userId, nil)

		userUCase := NewUseCase(authClient, userRepo)

		_, err := userUCase.Authorize(loginUser)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("Authorize_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		authClient := auth_service.NewMockAuthServiceClient(ctrl)

		authClient.
			EXPECT().
			Login(context.Background(), reqLoginUser).
			Return(userId, errors.ErrInternalError)

		userUCase := NewUseCase(authClient, userRepo)

		_, err := userUCase.Authorize(loginUser)
		assert.Error(t, err, "expected error")
	})
}

func TestUserUseCase_GetAvatar(t *testing.T) {
	profileUser := models.ProfileUser{}
	userId := uint64(323)

	t.Run("GetAvatar_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(userId).
			Return(&profileUser, nil)

		authClient := auth_service.NewMockAuthServiceClient(ctrl)

		userUCase := NewUseCase(authClient, userRepo)

		_, err := userUCase.GetAvatar(userId)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetAvatar_not_found_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(userId).
			Return(&profileUser, errors.ErrUserNotFound)

		authClient := auth_service.NewMockAuthServiceClient(ctrl)

		userUCase := NewUseCase(authClient, userRepo)

		_, err := userUCase.GetAvatar(userId)
		assert.Error(t, err, "expected error")
	})
}

func TestUserUseCase_GetUserById(t *testing.T) {
	profileUser := models.ProfileUser{}
	userId := uint64(323)

	t.Run("GetAvatar_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(userId).
			Return(&profileUser, nil)

		authClient := auth_service.NewMockAuthServiceClient(ctrl)

		userUCase := NewUseCase(authClient, userRepo)

		_, err := userUCase.GetUserById(userId)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetUserById_not_found_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(userId).
			Return(&profileUser, errors.ErrUserNotFound)

		authClient := auth_service.NewMockAuthServiceClient(ctrl)

		userUCase := NewUseCase(authClient, userRepo)

		_, err := userUCase.GetUserById(userId)
		assert.Error(t, err, "expected error")
	})
}

func TestUserUseCase_AddUser(t *testing.T) {
	userId := auth_service.UserId{Id: 0}
	profile := models.ProfileUser{Email: "test"}
	reqSignupUser := auth_service.SignupUser{
		Email:    "test",
		Password: "name",
	}
	signupUser := models.SignupUser{
		Email:    "test",
		Password: "name",
	}

	t.Run("AddUser_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			AddProfile(&profile).
			Return(uint64(0), nil)

		authClient := auth_service.NewMockAuthServiceClient(ctrl)
		authClient.
			EXPECT().
			Signup(context.Background(), &reqSignupUser).
			Return(&userId, nil)

		userUCase := NewUseCase(authClient, userRepo)

		_, err := userUCase.AddUser(&signupUser)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("AddUser_service_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)

		authClient := auth_service.NewMockAuthServiceClient(ctrl)
		authClient.
			EXPECT().
			Signup(context.Background(), &reqSignupUser).
			Return(&userId, errors.ErrInternalError)

		userUCase := NewUseCase(authClient, userRepo)

		_, err := userUCase.AddUser(&signupUser)
		assert.Error(t, err, "expected error")
	})

	t.Run("AddUser_add_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := user_mock.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			AddProfile(&profile).
			Return(uint64(0), errors.ErrInternalError)

		authClient := auth_service.NewMockAuthServiceClient(ctrl)
		authClient.
			EXPECT().
			Signup(context.Background(), &reqSignupUser).
			Return(&userId, nil)

		userUCase := NewUseCase(authClient, userRepo)

		_, err := userUCase.AddUser(&signupUser)
		assert.Error(t, err, "expected error")
	})
}
