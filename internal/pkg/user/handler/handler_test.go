package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
		Password: "qwerty",
	}
	invalidLoginUser := models.LoginUser{
		Email:    "test",
		Password: "qw",
	}
	profileUser := models.ProfileUser{
		Id:        3,
		FirstName: "test",
		LastName:  "last",
		Email:     "test@test.ru",
		Password:  []byte{1, 43, 23},
		Avatar: models.Avatar{
			Url: "http://test.png",
		},
	}
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: 3,
		},
	}

	t.Run("Login_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			Authorize(&loginUser).
			Return(&profileUser, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			CreateNewSession(profileUser.Id).
			Return(&sess, nil)

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesLogin, _ := json.Marshal(loginUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBuffer(bytesLogin))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Login)
		handler.ServeHTTP(rr, req)
	})

	t.Run("Login_incorrect_data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBuffer([]byte("test")))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("Login_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBufferString("ds"))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("Login_invalid_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesLogin, _ := json.Marshal(invalidLoginUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBuffer(bytesLogin))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("Login_incorrect_auth", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			Authorize(&loginUser).
			Return(nil, errors.ErrUserUnauthorized)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesLogin, _ := json.Marshal(loginUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBuffer(bytesLogin))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusUnauthorized, "incorrect http code")
	})

	t.Run("Login_server_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			Authorize(&loginUser).
			Return(&profileUser, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			CreateNewSession(profileUser.Id).
			Return(nil, errors.ErrInternalError)

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesLogin, _ := json.Marshal(loginUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/login",
			bytes.NewBuffer(bytesLogin))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	updateUser := models.UpdateUser{
		FirstName: "newName",
		LastName:  "newLast",
	}
	invalidUpdateUser := models.UpdateUser{
		FirstName: "ne",
		LastName:  "newLast",
	}
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: 3,
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

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesLogin, _ := json.Marshal(updateUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/user/profile",
			bytes.NewBuffer(bytesLogin))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.UpdateProfile)
		handler.ServeHTTP(rr, req)
	})

	t.Run("UpdateProfile_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/user/profile",
			bytes.NewBuffer([]byte("test")))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.UpdateProfile)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("UpdateProfile_invalid_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesLogin, _ := json.Marshal(invalidUpdateUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/user/profile",
			bytes.NewBuffer(bytesLogin))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.UpdateProfile)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("UpdateProfile_invalid_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			UpdateProfile(userId, &updateUser).
			Return(errors.ErrInternalError)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesLogin, _ := json.Marshal(updateUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/user/profile",
			bytes.NewBuffer(bytesLogin))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.UpdateProfile)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_GetProfileAvatar(t *testing.T) {
	avatar := models.Avatar{
		Url: "test.png",
	}
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("GetProfileAvatar_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			GetAvatar(userId).
			Return(avatar.Url, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/user/profile",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.GetProfileAvatar)
		handler.ServeHTTP(rr, req)
	})

	t.Run("GetProfileAvatar_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			GetAvatar(userId).
			Return("nil", errors.ErrInternalError)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/user/profile",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.GetProfileAvatar)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_GetProfile(t *testing.T) {
	userId := uint64(3)
	profileUser := models.ProfileUser{
		Id:        userId,
		FirstName: "test",
		LastName:  "last",
		Email:     "test@test.ru",
		Password:  []byte{1, 43, 23},
		Avatar: models.Avatar{
			Url: "http://test.png",
		},
	}
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("GetProfile_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			GetUserById(userId).
			Return(&profileUser, nil)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/profile",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.GetProfile)
		handler.ServeHTTP(rr, req)
	})

	t.Run("GetProfile_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			GetUserById(userId).
			Return(&profileUser, errors.ErrInternalError)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/profile",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.GetProfile)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_Signup(t *testing.T) {
	userId := uint64(3)
	signupUser := models.SignupUser{
		Email:    "test@test.ru",
		Password: "qwerty",
	}
	invalidSignupUser := models.SignupUser{
		Email:    "test.ru",
		Password: "qwerty",
	}
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

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesSignup, _ := json.Marshal(signupUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/signup",
			bytes.NewBuffer(bytesSignup))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Signup)
		handler.ServeHTTP(rr, req)
	})

	t.Run("Signup_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/signup",
			bytes.NewBuffer([]byte("test")))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Signup)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("Signup_invalid_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesLogin, _ := json.Marshal(invalidSignupUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/signup",
			bytes.NewBuffer(bytesLogin))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Signup)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("Signup_conflict_adding", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)
		userUCase.
			EXPECT().
			AddUser(&signupUser).
			Return(userId, errors.ErrInternalError)

		sessionUCase := session_mock.NewMockUseCase(ctrl)

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesSignup, _ := json.Marshal(signupUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/signup",
			bytes.NewBuffer(bytesSignup))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Signup)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusConflict, "incorrect http code")
	})

	t.Run("Signup_server_error", func(t *testing.T) {
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

		userHandler := NewHandler(userUCase, sessionUCase)

		bytesSignup, _ := json.Marshal(signupUser)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/user/signup",
			bytes.NewBuffer(bytesSignup))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Signup)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_Logout(t *testing.T) {
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("Logout_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)

		sessionUCase := session_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			DestroySession(sess.Value).
			Return(nil)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/user/logout",
			bytes.NewBuffer(nil))
		req.AddCookie(&http.Cookie{
			Name:    models.SessionCookieName,
			Value:   sess.Value,
			Expires: time.Now().Add(models.ExpireSessionCookie * time.Second),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Logout)
		handler.ServeHTTP(rr, req)
	})

	t.Run("Logout_user_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userUCase := user_mock.NewMockUseCase(ctrl)

		sessionUCase := session_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			DestroySession(sess.Value).
			Return(errors.ErrInternalError)

		userHandler := NewHandler(userUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/user/logout",
			bytes.NewBuffer(nil))
		req.AddCookie(&http.Cookie{
			Name:    models.SessionCookieName,
			Value:   sess.Value,
			Expires: time.Now().Add(models.ExpireSessionCookie * time.Second),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(userHandler.Logout)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}
