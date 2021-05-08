package product

import (
	category_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/mock"
	product_repo "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/mock"
	"testing"

	category_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	product_mock "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product/mock"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProductUseCase_GetProductById(t *testing.T) {
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

func TestProductUseCase_GetRangeProducts(t *testing.T) {
	paginator := models.PaginatorProducts{
		PageNum:     1,
		Count:       12,
		Category:    4,
		Filter:      &models.ProductFilter{
			MinPrice:   0,
			MaxPrice:   10,
			IsNew:      false,
			IsRating:   false,
			IsDiscount: false,
		},
		SortOptions: models.SortOptions{
			SortKey:       "date",
			SortDirection: "ASC",
		},
	}
	products := []*models.ViewProduct{{}}

	t.Run("GetRangeProducts_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_repo.NewMockRepository(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)

		productRepo.
			EXPECT().
			CreateFilterString(paginator.Filter).
			Return("")

		productRepo.
			EXPECT().
			GetCountPages(paginator.Category, paginator.Count, "").
			Return(10, nil)

		productRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", nil)

		productRepo.
			EXPECT().
			SelectRangeProducts(&paginator, "", "").
			Return(products, nil)

		productUCase := NewUseCase(productRepo, categoryRepo)

		_, err := productUCase.GetRangeProducts(&paginator)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetRangeProducts_bad_count", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_repo.NewMockRepository(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)

		productRepo.
			EXPECT().
			CreateFilterString(paginator.Filter).
			Return("")

		productRepo.
			EXPECT().
			GetCountPages(paginator.Category, paginator.Count, "").
			Return(10, errors.ErrInternalError)

		productUCase := NewUseCase(productRepo, categoryRepo)

		_, err := productUCase.GetRangeProducts(&paginator)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeProducts_bad_sort_options", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_repo.NewMockRepository(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)

		productRepo.
			EXPECT().
			CreateFilterString(paginator.Filter).
			Return("")

		productRepo.
			EXPECT().
			GetCountPages(paginator.Category, paginator.Count, "").
			Return(10, nil)

		productRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", errors.ErrInternalError)

		productUCase := NewUseCase(productRepo, categoryRepo)

		_, err := productUCase.GetRangeProducts(&paginator)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeProducts_not_found_products", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_repo.NewMockRepository(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)

		productRepo.
			EXPECT().
			CreateFilterString(paginator.Filter).
			Return("")

		productRepo.
			EXPECT().
			GetCountPages(paginator.Category, paginator.Count, "").
			Return(10, nil)

		productRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", nil)

		productRepo.
			EXPECT().
			SelectRangeProducts(&paginator, "", "").
			Return(products, errors.ErrInternalError)

		productUCase := NewUseCase(productRepo, categoryRepo)

		_, err := productUCase.GetRangeProducts(&paginator)
		assert.Error(t, err, "expected error")
	})
}

func TestProductUseCase_SearchRangeProducts(t *testing.T) {
	paginator := models.SearchQuery{
		QueryString: "test",
		PageNum:     1,
		Count:       12,
		Category:    4,
		Filter:      &models.ProductFilter{
			MinPrice:   0,
			MaxPrice:   10,
			IsNew:      false,
			IsRating:   false,
			IsDiscount: false,
		},
		SortOptions: models.SortOptions{
			SortKey:       "date",
			SortDirection: "ASC",
		},
	}
	products := []*models.ViewProduct{{}}

	t.Run("SearchRangeProducts_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_repo.NewMockRepository(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)

		productRepo.
			EXPECT().
			CreateFilterString(paginator.Filter).
			Return("")

		productRepo.
			EXPECT().
			GetCountSearchPages(paginator.Category, paginator.Count,
				paginator.QueryString, "").
			Return(10, nil)

		productRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", nil)

		productRepo.
			EXPECT().
			SearchRangeProducts(&paginator, "", "").
			Return(products, nil)

		productUCase := NewUseCase(productRepo, categoryRepo)

		_, err := productUCase.SearchRangeProducts(&paginator)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetRangeProducts_bad_count", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_repo.NewMockRepository(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)

		productRepo.
			EXPECT().
			CreateFilterString(paginator.Filter).
			Return("")

		productRepo.
			EXPECT().
			GetCountSearchPages(paginator.Category, paginator.Count,
				paginator.QueryString, "").
			Return(10, errors.ErrInternalError)

		productUCase := NewUseCase(productRepo, categoryRepo)

		_, err := productUCase.SearchRangeProducts(&paginator)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeProducts_bad_sort_options", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_repo.NewMockRepository(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)

		productRepo.
			EXPECT().
			CreateFilterString(paginator.Filter).
			Return("")

		productRepo.
			EXPECT().
			GetCountSearchPages(paginator.Category, paginator.Count,
				paginator.QueryString, "").
			Return(10, nil)

		productRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", errors.ErrInternalError)

		productUCase := NewUseCase(productRepo, categoryRepo)

		_, err := productUCase.SearchRangeProducts(&paginator)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeProducts_not_found_products", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_repo.NewMockRepository(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)

		productRepo.
			EXPECT().
			CreateFilterString(paginator.Filter).
			Return("")

		productRepo.
			EXPECT().
			GetCountSearchPages(paginator.Category, paginator.Count,
				paginator.QueryString, "").
			Return(10, nil)

		productRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", nil)

		productRepo.
			EXPECT().
			SearchRangeProducts(&paginator, "", "").
			Return(products, errors.ErrInternalError)

		productUCase := NewUseCase(productRepo, categoryRepo)

		_, err := productUCase.SearchRangeProducts(&paginator)
		assert.Error(t, err, "expected error")
	})
}