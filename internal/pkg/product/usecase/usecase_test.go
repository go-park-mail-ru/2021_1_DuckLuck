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
