package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cart_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestCartHandler_GetProduct(t *testing.T) {
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}
	cartArticle := models.CartArticle{
		ProductPosition: models.ProductPosition{
			Count: 3,
		},
		ProductIdentifier: models.ProductIdentifier{
			ProductId: 2,
		},
	}

	t.Run("AddProductInCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			AddProduct(sess.UserData.Id, &cartArticle).
			Return(nil)

		cartHandler := NewHandler(cartUCase)

		bytesArticle, _ := json.Marshal(cartArticle)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/cart/product",
			bytes.NewBuffer(bytesArticle))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.AddProductInCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("AddProductInCart_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)

		cartHandler := NewHandler(cartUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/cart/product",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.AddProductInCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("AddProductInCart_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			AddProduct(sess.UserData.Id, &cartArticle).
			Return(errors.ErrInternalError)

		cartHandler := NewHandler(cartUCase)

		bytesArticle, _ := json.Marshal(cartArticle)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/cart/product",
			bytes.NewBuffer(bytesArticle))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.AddProductInCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestCartHandler_DeleteProductInCart(t *testing.T) {
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}
	identifier := models.ProductIdentifier{
		ProductId: 3,
	}

	t.Run("DeleteProductInCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			DeleteProduct(sess.UserData.Id, &identifier).
			Return(nil)

		cartHandler := NewHandler(cartUCase)

		bytesIdentifier, _ := json.Marshal(identifier)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/cart/product",
			bytes.NewBuffer(bytesIdentifier))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.DeleteProductInCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("DeleteProductInCart_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)

		cartHandler := NewHandler(cartUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/cart/product",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.DeleteProductInCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("DeleteProductInCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			DeleteProduct(sess.UserData.Id, &identifier).
			Return(errors.ErrInternalError)

		cartHandler := NewHandler(cartUCase)

		bytesIdentifier, _ := json.Marshal(identifier)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/cart/product",
			bytes.NewBuffer(bytesIdentifier))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.DeleteProductInCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestCartHandler_ChangeProductInCart(t *testing.T) {
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}
	cartArticle := models.CartArticle{
		ProductPosition: models.ProductPosition{
			Count: 3,
		},
		ProductIdentifier: models.ProductIdentifier{
			ProductId: 2,
		},
	}

	t.Run("ChangeProductInCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			ChangeProduct(sess.UserData.Id, &cartArticle).
			Return(nil)

		cartHandler := NewHandler(cartUCase)

		bytesArticle, _ := json.Marshal(cartArticle)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/cart/product",
			bytes.NewBuffer(bytesArticle))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.ChangeProductInCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("ChangeProductInCart_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)

		cartHandler := NewHandler(cartUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/cart/product",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.ChangeProductInCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("ChangeProductInCart_can't_change_product", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			ChangeProduct(sess.UserData.Id, &cartArticle).
			Return(errors.ErrInternalError)

		cartHandler := NewHandler(cartUCase)

		bytesArticle, _ := json.Marshal(cartArticle)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "PUT", "/api/v1/cart/product",
			bytes.NewBuffer(bytesArticle))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.ChangeProductInCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestCartHandler_GetProductsFromCart(t *testing.T) {
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}
	previewCart := models.PreviewCart{
		Products: nil,
		Price: models.TotalPrice{
			TotalDiscount: 12,
			TotalCost:     23,
			TotalBaseCost: 43,
		},
	}

	t.Run("GetProductsFromCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			GetPreviewCart(sess.UserData.Id).
			Return(&previewCart, nil)

		cartHandler := NewHandler(cartUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/cart",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.GetProductsFromCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("GetProductsFromCart_can't_found_cart", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			GetPreviewCart(sess.UserData.Id).
			Return(&previewCart, errors.ErrInternalError)

		cartHandler := NewHandler(cartUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/cart",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.GetProductsFromCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestCartHandler_DeleteProductsFromCart(t *testing.T) {
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("DeleteProductsFromCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			DeleteCart(sess.UserData.Id).
			Return(nil)

		cartHandler := NewHandler(cartUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/cart",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.DeleteProductsFromCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("DeleteProductsFromCart_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cartUCase := cart_mock.NewMockUseCase(ctrl)
		cartUCase.
			EXPECT().
			DeleteCart(sess.UserData.Id).
			Return(errors.ErrInternalError)

		cartHandler := NewHandler(cartUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/cart",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(cartHandler.DeleteProductsFromCart)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}
