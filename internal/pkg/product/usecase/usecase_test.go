package product

import (
	"testing"

	category_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	product_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserUseCase_GetProductById(t *testing.T) {
	product := models.Product{
		Id:    2,
		Title: "test",
		Price: models.ProductPrice{
			Discount: 10,
			BaseCost: 20,
		},
		Rating:       5,
		Description:  "fdfdfddf",
		Category:     4,
		CategoryPath: nil,
		Images:       nil,
	}
	categories := []*models.CategoriesCatalog{
		&models.CategoriesCatalog{
			Id:   3,
			Name: "home",
			Next: nil,
		},
		&models.CategoriesCatalog{
			Id:   3,
			Name: "electronic",
			Next: nil,
		},
	}
	resProduct := models.Product{
		Id:    2,
		Title: "test",
		Price: models.ProductPrice{
			Discount: 10,
			BaseCost: 20,
		},
		Rating:       5,
		Description:  "fdfdfddf",
		Category:     4,
		CategoryPath: categories,
		Images:       nil,
	}

	t.Run("GetUserIdBySession_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		productRepo.
			EXPECT().
			SelectProductById(product.Id).
			Return(&product, nil)

		categoryRepo := category_mock.NewMockRepository(ctrl)
		categoryRepo.
			EXPECT().
			GetPathToCategory(product.Category).
			Return(categories, nil)

		userUCase := NewUseCase(productRepo, categoryRepo)

		userData, err := userUCase.GetProductById(product.Id)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, resProduct, *userData, "not equal data")
	})

	t.Run("GetUserIdBySession_product_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		productRepo.
			EXPECT().
			SelectProductById(product.Id).
			Return(&product, errors.ErrDBInternalError)

		categoryRepo := category_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(productRepo, categoryRepo)

		_, err := userUCase.GetProductById(product.Id)
		assert.Equal(t, err, errors.ErrProductNotFound, "not equal errors")
	})

	t.Run("GetUserIdBySession_category_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productRepo := product_mock.NewMockRepository(ctrl)
		productRepo.
			EXPECT().
			SelectProductById(product.Id).
			Return(&product, nil)

		categoryRepo := category_mock.NewMockRepository(ctrl)
		categoryRepo.
			EXPECT().
			GetPathToCategory(product.Category).
			Return(categories, errors.ErrDBInternalError)

		userUCase := NewUseCase(productRepo, categoryRepo)

		_, err := userUCase.GetProductById(product.Id)
		assert.Equal(t, err, errors.ErrCategoryNotFound, "not equal errors")
	})
}

func TestUserUseCase_GetRangeProducts(t *testing.T) {
	paginator := models.PaginatorProducts{
		PageNum:       1,
		Count:         10,
		SortKey:       "cost",
		SortDirection: "ASC",
		Category:      1,
	}
	incorrectPaginator := models.PaginatorProducts{
		PageNum:       0,
		Count:         0,
		SortKey:       "cost",
		SortDirection: "ASC",
		Category:      1,
	}
	categories := []uint64{1, 2, 4, 5}
	rangeProduct := models.RangeProducts{
		ListPreviewProducts: []*models.ViewProduct{
			&models.ViewProduct{
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

	t.Run("GetRangeProducts_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_mock.NewMockRepository(ctrl)
		categoryRepo.
			EXPECT().
			GetAllSubCategoriesId(paginator.Category).
			Return(categories, nil)

		productRepo := product_mock.NewMockRepository(ctrl)
		productRepo.
			EXPECT().
			SelectRangeProducts(&paginator, &categories).
			Return(&rangeProduct, nil)

		userUCase := NewUseCase(productRepo, categoryRepo)

		userData, err := userUCase.GetRangeProducts(&paginator)
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, rangeProduct, *userData, "not equal data")
	})

	t.Run("GetRangeProducts_incorrect_paginator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_mock.NewMockRepository(ctrl)
		productRepo := product_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(productRepo, categoryRepo)

		_, err := userUCase.GetRangeProducts(&incorrectPaginator)
		assert.Equal(t, err, errors.ErrIncorrectPaginator, "not equal data")
	})

	t.Run("GetRangeProducts_product_not_found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_mock.NewMockRepository(ctrl)
		categoryRepo.
			EXPECT().
			GetAllSubCategoriesId(paginator.Category).
			Return(nil, errors.ErrDBInternalError)

		productRepo := product_mock.NewMockRepository(ctrl)

		userUCase := NewUseCase(productRepo, categoryRepo)

		_, err := userUCase.GetRangeProducts(&paginator)
		assert.Equal(t, err, errors.ErrCategoryNotFound, "not equal data")
	})
}
