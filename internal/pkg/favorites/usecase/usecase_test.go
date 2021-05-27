package usecase

import (
	favorites_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/favorites/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFavoritesUseCase_AddProductToFavorites(t *testing.T) {
	productId := uint64(1)
	userId := uint64(3)

	t.Run("AddProductToFavorites_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesRepo := favorites_mock.NewMockRepository(ctrl)
		favoritesRepo.
			EXPECT().
			AddProductToFavorites(productId, userId).
			Return(nil)

		favoritesUCase := NewUseCase(favoritesRepo)

		err := favoritesUCase.AddProductToFavorites(productId, userId)
		assert.NoError(t, err, "unexpected error")
	})
}

func TestFavoritesUseCase_DeleteProductFromFavorites(t *testing.T) {
	productId := uint64(1)
	userId := uint64(3)

	t.Run("DeleteProductFromFavorites_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesRepo := favorites_mock.NewMockRepository(ctrl)
		favoritesRepo.
			EXPECT().
			DeleteProductFromFavorites(productId, userId).
			Return(nil)

		favoritesUCase := NewUseCase(favoritesRepo)

		err := favoritesUCase.DeleteProductFromFavorites(productId, userId)
		assert.NoError(t, err, "unexpected error")
	})
}

func TestFavoritesUseCase_GetUserFavorites(t *testing.T) {
	userId := uint64(3)

	t.Run("GetUserFavorites_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesRepo := favorites_mock.NewMockRepository(ctrl)
		favoritesRepo.
			EXPECT().
			GetUserFavorites(userId).
			Return(nil, nil)

		favoritesUCase := NewUseCase(favoritesRepo)

		_, err := favoritesUCase.GetUserFavorites(userId)
		assert.NoError(t, err, "unexpected error")
	})
}

func TestFavoritesUseCase_GetRangeFavorites(t *testing.T) {
	userId := uint64(3)
	paginator := models.PaginatorFavorites{
		PageNum:     2,
		Count:       3,
		SortOptions: models.SortOptions{},
	}
	badPaginator := models.PaginatorFavorites{}
	countPages := 10
	sortString := ""

	t.Run("GetRangeFavorites_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesRepo := favorites_mock.NewMockRepository(ctrl)

		favoritesRepo.
			EXPECT().
			GetCountPages(userId, paginator.Count).
			Return(countPages, nil)

		favoritesRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return(sortString, nil)

		favoritesRepo.
			EXPECT().
			SelectRangeFavorites(&paginator, paginator.SortKey, userId).
			Return(nil, nil)

		favoritesUCase := NewUseCase(favoritesRepo)

		_, err := favoritesUCase.GetRangeFavorites(&paginator, userId)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetRangeFavorites_incorrect_paginator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesRepo := favorites_mock.NewMockRepository(ctrl)
		favoritesUCase := NewUseCase(favoritesRepo)

		_, err := favoritesUCase.GetRangeFavorites(&badPaginator, userId)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeFavorites_get_count_pages_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesRepo := favorites_mock.NewMockRepository(ctrl)

		favoritesRepo.
			EXPECT().
			GetCountPages(userId, paginator.Count).
			Return(countPages, errors.ErrInternalError)

		favoritesUCase := NewUseCase(favoritesRepo)

		_, err := favoritesUCase.GetRangeFavorites(&paginator, userId)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeFavorites_can't_create_sort_string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesRepo := favorites_mock.NewMockRepository(ctrl)

		favoritesRepo.
			EXPECT().
			GetCountPages(userId, paginator.Count).
			Return(countPages, nil)

		favoritesRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return(sortString, errors.ErrDBInternalError)

		favoritesUCase := NewUseCase(favoritesRepo)

		_, err := favoritesUCase.GetRangeFavorites(&paginator, userId)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeFavorites_can't_select_range", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		favoritesRepo := favorites_mock.NewMockRepository(ctrl)

		favoritesRepo.
			EXPECT().
			GetCountPages(userId, paginator.Count).
			Return(countPages, nil)

		favoritesRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return(sortString, nil)

		favoritesRepo.
			EXPECT().
			SelectRangeFavorites(&paginator, paginator.SortKey, userId).
			Return(nil, errors.ErrDBInternalError)

		favoritesUCase := NewUseCase(favoritesRepo)

		_, err := favoritesUCase.GetRangeFavorites(&paginator, userId)
		assert.Error(t, err, "expected error")
	})
}
