package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	notification_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/notification/mock"

	"github.com/golang/mock/gomock"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestFavoritesHandler_SubscribeUser(t *testing.T) {
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}
	notificationCredentials := models.NotificationCredentials{
		UserIdentifier: models.UserIdentifier{},
		Keys:           models.NotificationKeys{},
	}

	t.Run("SubscribeUser_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationUCase := notification_mock.NewMockUseCase(ctrl)
		notificationUCase.
			EXPECT().
			SubscribeUser(userId, &notificationCredentials).
			Return(nil)

		notificationHandler := NewHandler(notificationUCase)

		bytesNotification, _ := json.Marshal(notificationCredentials)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(bytesNotification))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(notificationHandler.SubscribeUser)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("SubscribeUser_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationUCase := notification_mock.NewMockUseCase(ctrl)
		notificationHandler := NewHandler(notificationUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(notificationHandler.SubscribeUser)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("SubscribeUser_can't_subscribe", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationUCase := notification_mock.NewMockUseCase(ctrl)
		notificationUCase.
			EXPECT().
			SubscribeUser(userId, &notificationCredentials).
			Return(errors.ErrInternalError)

		notificationHandler := NewHandler(notificationUCase)

		bytesNotification, _ := json.Marshal(notificationCredentials)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(bytesNotification))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(notificationHandler.SubscribeUser)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestFavoritesHandler_UnsubscribeUser(t *testing.T) {
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}
	userIdentifier := models.UserIdentifier{}

	t.Run("UnsubscribeUser_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationUCase := notification_mock.NewMockUseCase(ctrl)
		notificationUCase.
			EXPECT().
			UnsubscribeUser(userId, userIdentifier.Endpoint).
			Return(nil)

		notificationHandler := NewHandler(notificationUCase)

		bytesNotification, _ := json.Marshal(userIdentifier)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(bytesNotification))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(notificationHandler.UnsubscribeUser)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("UnsubscribeUser_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationUCase := notification_mock.NewMockUseCase(ctrl)
		notificationHandler := NewHandler(notificationUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(notificationHandler.UnsubscribeUser)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("UnsubscribeUser_can't_subscribe", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationUCase := notification_mock.NewMockUseCase(ctrl)
		notificationUCase.
			EXPECT().
			UnsubscribeUser(userId, userIdentifier.Endpoint).
			Return(errors.ErrInternalError)

		notificationHandler := NewHandler(notificationUCase)

		bytesNotification, _ := json.Marshal(userIdentifier)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(bytesNotification))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(notificationHandler.UnsubscribeUser)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestFavoritesHandler_GetNotificationPublicKey(t *testing.T) {
	t.Run("GetNotificationPublicKey_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		notificationUCase := notification_mock.NewMockUseCase(ctrl)
		notificationHandler := NewHandler(notificationUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/notification/key",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(notificationHandler.GetNotificationPublicKey)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})
}
