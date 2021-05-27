package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	favorites_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/favorites/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestFavoritesHandler_AddProductToFavorites(t *testing.T) {
	productId := uint64(1)
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("AddProductToFavorites_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesUCase.
			EXPECT().
			AddProductToFavorites(productId, userId).
			Return(nil)

		favoritesHandler := NewHandler(favoritesUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.AddProductToFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("AddProductToFavorites_bad_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesHandler := NewHandler(favoritesUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"_id_": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.AddProductToFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("AddProductToFavorites_product_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesUCase.
			EXPECT().
			AddProductToFavorites(productId, userId).
			Return(errors.ErrProductNotFound)

		favoritesHandler := NewHandler(favoritesUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.AddProductToFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestFavoritesHandler_DeleteProductFromFavorites(t *testing.T) {
	productId := uint64(1)
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("DeleteProductFromFavorites_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesUCase.
			EXPECT().
			DeleteProductFromFavorites(productId, userId).
			Return(nil)

		favoritesHandler := NewHandler(favoritesUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.DeleteProductFromFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("DeleteProductFromFavorites_bad_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesHandler := NewHandler(favoritesUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"_id_": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.DeleteProductFromFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("DeleteProductFromFavorites_product_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesUCase.
			EXPECT().
			DeleteProductFromFavorites(productId, userId).
			Return(errors.ErrProductNotFound)

		favoritesHandler := NewHandler(favoritesUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "DELETE", "/api/v1/favorites/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.DeleteProductFromFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestFavoritesHandler_GetListPreviewFavorites(t *testing.T) {
	paginator := models.PaginatorFavorites{
		PageNum: 2,
		Count:   10,
		SortOptions: models.SortOptions{
			SortKey:       "cost",
			SortDirection: "ASC",
		},
	}
	badPaginator := models.PaginatorFavorites{
		PageNum: 2,
		Count:   10,
		SortOptions: models.SortOptions{
			SortKey:       "cost",
			SortDirection: "TEST",
		},
	}
	rangeFavorites := models.RangeFavorites{
		ListPreviewProducts: []*models.ViewFavorite{
			&models.ViewFavorite{
				Id:    3,
				Title: "test",
				Price: models.ProductPrice{
					Discount: 10,
					BaseCost: 50,
				},
				Rating:       5,
				PreviewImage: "fdfdf",
			},
		},
		MaxCountPages: 3,
	}
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}

	t.Run("GetListPreviewFavorites_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesUCase.
			EXPECT().
			GetRangeFavorites(&paginator, userId).
			Return(&rangeFavorites, nil)

		favoritesHandler := NewHandler(favoritesUCase)

		bytesPaginator, _ := json.Marshal(paginator)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites",
			bytes.NewBuffer(bytesPaginator))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.GetListPreviewFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("GetListPreviewFavorites_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesHandler := NewHandler(favoritesUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.GetListPreviewFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("GetListPreviewFavorites_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesUCase.
			EXPECT().
			GetRangeFavorites(&paginator, userId).
			Return(&rangeFavorites, errors.ErrInternalError)

		favoritesHandler := NewHandler(favoritesUCase)

		bytesPaginator, _ := json.Marshal(paginator)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites",
			bytes.NewBuffer(bytesPaginator))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.GetListPreviewFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})

	t.Run("GetListPreviewFavorites_bad_paginator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesHandler := NewHandler(favoritesUCase)

		bytesPaginator, _ := json.Marshal(badPaginator)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "POST", "/api/v1/favorites",
			bytes.NewBuffer(bytesPaginator))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.GetListPreviewFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})
}

func TestFavoritesHandler_GetUserFavorites(t *testing.T) {
	userId := uint64(3)
	sess := models.Session{
		Value: "fdsfdsfdsf",
		UserData: models.UserId{
			Id: userId,
		},
	}
	userFavorites := models.UserFavorites{Products: nil}

	t.Run("GetUserFavorites_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesUCase.
			EXPECT().
			GetUserFavorites(userId).
			Return(&userFavorites, nil)

		favoritesHandler := NewHandler(favoritesUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/favorites",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.GetUserFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("GetUserFavorites_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesUCase := favorites_mock.NewMockUseCase(ctrl)
		favoritesUCase.
			EXPECT().
			GetUserFavorites(userId).
			Return(nil, errors.ErrProductNotFound)

		favoritesHandler := NewHandler(favoritesUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		cctx := context.WithValue(ctx, models.SessionContextKey, &sess)
		req, _ := http.NewRequestWithContext(cctx, "GET", "/api/v1/favorites",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(favoritesHandler.GetUserFavorites)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}
