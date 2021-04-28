package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	cart_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	order_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/order/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestOrderHandler_GetOrderFromCart(t *testing.T) {
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: 3,
		},
	}
	previewCart := models.PreviewCart{
		Products: nil,
		Price: models.TotalPrice{
			TotalDiscount: 32,
			TotalCost:     234,
			TotalBaseCost: 34,
		},
	}
	previewOrder := models.PreviewOrder{
		Products: nil,
		Recipient: models.OrderRecipient{
			FirstName: "test",
			LastName:  "last",
			Email:     "test@test.ru",
		},
		Price: models.TotalPrice{
			TotalDiscount: 32,
			TotalCost:     234,
			TotalBaseCost: 34,
		},
		Address: models.OrderAddress{
			Address: "",
		},
	}

	t.Run("GetOrderFromCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			GetPreviewCart(sess.UserData.Id).
			Return(&previewCart, nil)

		orderUCase := order_mock.NewMockUseCase(ctrl)
		orderUCase.
			EXPECT().
			GetPreviewOrder(sess.UserData.Id, &previewCart).
			Return(&previewOrder, nil)

		orderHandler := NewHandler(orderUCase, cartUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/order",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(orderHandler.GetOrderFromCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("GetOrderFromCart_not_found_preview_cart", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionUCase := cart_mock.NewMockUseCase(ctrl)
		sessionUCase.
			EXPECT().
			GetPreviewCart(sess.UserData.Id).
			Return(&previewCart, nil)

		orderUCase := order_mock.NewMockUseCase(ctrl)
		orderUCase.
			EXPECT().
			GetPreviewOrder(sess.UserData.Id, &previewCart).
			Return(&previewOrder, errors.ErrInternalError)

		orderHandler := NewHandler(orderUCase, sessionUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/order",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(orderHandler.GetOrderFromCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})

	t.Run("GetOrderFromCart_not_found_preview_cart", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			GetPreviewCart(sess.UserData.Id).
			Return(&previewCart, errors.ErrInternalError)

		orderUCase := order_mock.NewMockUseCase(ctrl)

		orderHandler := NewHandler(orderUCase, cartUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/order",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(orderHandler.GetOrderFromCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}
