package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	session_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/session/mock"
	user_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/user/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Login(t *testing.T) {
	loginUser := models.LoginUser{
		Email:    "test@test.ru",
		Password: "fvdfvdf",
	}
	userId := uint64(12)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("Login_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			Authorize(&loginUser).
			Return(userId, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			CreateNewSession(userId).
			Return(&sess, nil)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		userBytes, _ := json.Marshal(loginUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBuffer(userBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("Login_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("Login_not_authorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			Authorize(&loginUser).
			Return(userId, errors.ErrInternalError)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		userBytes, _ := json.Marshal(loginUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBuffer(userBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusUnauthorized, "incorrect http code")
	})

	t.Run("Login_not_create_session", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			Authorize(&loginUser).
			Return(userId, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			CreateNewSession(userId).
			Return(&sess, errors.ErrInternalError)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		userBytes, _ := json.Marshal(loginUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBuffer(userBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	updateUser := models.UpdateUser{
		FirstName: "test",
		LastName:  "name",
	}
	userId := uint64(12)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("UpdateProfile_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			UpdateProfile(userId, &updateUser).
			Return(nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		userBytes, _ := json.Marshal(updateUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/user/profile",
			bytes.NewBuffer(userBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.UpdateProfile)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("UpdateProfile_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/user/profile",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.UpdateProfile)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("UpdateProfile_not_update_profile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			UpdateProfile(userId, &updateUser).
			Return(errors.ErrInternalError)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		userBytes, _ := json.Marshal(updateUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/user/profile",
			bytes.NewBuffer(userBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.UpdateProfile)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_UpdateProfileAvatar(t *testing.T) {
	userId := uint64(12)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("UpdateProfileAvatar_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/user/profile/avatar",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.UpdateProfileAvatar)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})
}

func TestUserHandler_GetProfileAvatar(t *testing.T) {
	userId := uint64(12)
	profileAvatar := "test_url"
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("UpdateProfile_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			GetAvatar(userId).
			Return(profileAvatar, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/user/profile/avatar",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.GetProfileAvatar)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("UpdateProfile_not_found_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			GetAvatar(userId).
			Return(profileAvatar, errors.ErrInternalError)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/user/profile/avatar",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.GetProfileAvatar)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_GetProfile(t *testing.T) {
	userId := uint64(12)
	profileUser := models.ProfileUser{
		Id:     userId,
		AuthId: 2323,
		Email:  "test@test.ru",
	}
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("UpdateProfile_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			GetUserById(userId).
			Return(&profileUser, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/user/profile",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.GetProfile)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("UpdateProfile_not_found_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			GetUserById(userId).
			Return(&profileUser, errors.ErrInternalError)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/user/profile",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.GetProfile)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_Signup(t *testing.T) {
	signupUser := models.SignupUser{
		Email:    "test@test.ru",
		Password: "fvdfvdf",
	}
	userId := uint64(12)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("Signup_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			AddUser(&signupUser).
			Return(userId, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			CreateNewSession(userId).
			Return(&sess, nil)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		userBytes, _ := json.Marshal(signupUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/signup",
			bytes.NewBuffer(userBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.Signup)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusCreated, "incorrect http code")
	})

	t.Run("Signup_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/signup",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.Signup)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("Signup_not_authorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			AddUser(&signupUser).
			Return(userId, errors.ErrInternalError)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		userBytes, _ := json.Marshal(signupUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/signup",
			bytes.NewBuffer(userBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.Signup)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusConflict, "incorrect http code")
	})

	t.Run("Signup_not_create_session", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			AddUser(&signupUser).
			Return(userId, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			CreateNewSession(userId).
			Return(&sess, errors.ErrInternalError)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		userBytes, _ := json.Marshal(signupUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/signup",
			bytes.NewBuffer(userBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.Signup)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_Logout(t *testing.T) {
	userId := uint64(12)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("Logout_session_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)

		sessionUCase := session_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			DestroySession(sess.Value).
			Return(errors.ErrInternalError)

		reviewHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/user/logout",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(reviewHandler.Logout)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}
