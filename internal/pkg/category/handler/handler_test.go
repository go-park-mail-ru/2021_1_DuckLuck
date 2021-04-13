package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	category_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetCatalogCategories(t *testing.T) {
	categories := []*models.CategoriesCatalog{
		&models.CategoriesCatalog{
			Id:   4,
			Name: "test",
			Next: nil,
		},
	}

	t.Run("GetCatalogCategories_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryUCase := category_mock.NewMockUseCase(ctrl)
		categoryUCase.
			EXPECT().
			GetCatalogCategories().
			Return(categories, nil)

		categoryHandler := NewHandler(categoryUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/category",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(categoryHandler.GetCatalogCategories)
		handler.ServeHTTP(rr, req)
	})

	t.Run("GetCatalogCategories_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryUCase := category_mock.NewMockUseCase(ctrl)
		categoryUCase.
			EXPECT().
			GetCatalogCategories().
			Return(categories, errors.ErrInternalError)

		categoryHandler := NewHandler(categoryUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/category",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(categoryHandler.GetCatalogCategories)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestUserHandler_GetSubCategories(t *testing.T) {
	categories := []*models.CategoriesCatalog{
		&models.CategoriesCatalog{
			Id:   4,
			Name: "test",
			Next: nil,
		},
	}
	categoryId := uint64(3)

	t.Run("GetSubCategories_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryUCase := category_mock.NewMockUseCase(ctrl)
		categoryUCase.
			EXPECT().
			GetSubCategoriesById(categoryId).
			Return(categories, nil)

		categoryHandler := NewHandler(categoryUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/category/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", categoryId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(categoryHandler.GetSubCategories)
		handler.ServeHTTP(rr, req)
	})

	t.Run("GetSubCategories_without_args", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryUCase := category_mock.NewMockUseCase(ctrl)

		categoryHandler := NewHandler(categoryUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/category/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(categoryHandler.GetSubCategories)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("GetSubCategories_categories_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryUCase := category_mock.NewMockUseCase(ctrl)
		categoryUCase.
			EXPECT().
			GetSubCategoriesById(categoryId).
			Return(nil, errors.ErrInternalError)

		categoryHandler := NewHandler(categoryUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/category/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", categoryId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(categoryHandler.GetSubCategories)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}
