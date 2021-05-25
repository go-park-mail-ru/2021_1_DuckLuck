package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	admin_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/admin/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

	"github.com/golang/mock/gomock"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestFavoritesHandler_ChangeOrderStatus(t *testing.T) {
	updateOrder := models.UpdateOrder{
		OrderId: 3,
		Status:  "получено",
	}

	t.Run("ChangeOrderStatus_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		adminUCase := admin_mock.NewMockUseCase(ctrl)
		adminUCase.
			EXPECT().
			ChangeOrderStatus(&updateOrder).
			Return(nil)

		adminHandler := NewHandler(adminUCase)

		orderBytes, _ := json.Marshal(updateOrder)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/admin/order/status",
			bytes.NewBuffer(orderBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(adminHandler.ChangeOrderStatus)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("ChangeOrderStatus_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		adminUCase := admin_mock.NewMockUseCase(ctrl)
		adminHandler := NewHandler(adminUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/admin/order/status",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(adminHandler.ChangeOrderStatus)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("ChangeOrderStatus_can't_change_status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		adminUCase := admin_mock.NewMockUseCase(ctrl)
		adminUCase.
			EXPECT().
			ChangeOrderStatus(&updateOrder).
			Return(errors.ErrInternalError)

		adminHandler := NewHandler(adminUCase)

		orderBytes, _ := json.Marshal(updateOrder)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/admin/order/status",
			bytes.NewBuffer(orderBytes))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(adminHandler.ChangeOrderStatus)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}
