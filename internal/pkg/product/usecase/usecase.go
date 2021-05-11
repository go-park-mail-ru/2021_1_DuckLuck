package product

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type ProductUseCase struct {
	ProductRepo  product.Repository
	CategoryRepo category.Repository
}

func NewUseCase(productRepo product.Repository, categoryRepo category.Repository) product.UseCase {
	return &ProductUseCase{
		ProductRepo:  productRepo,
		CategoryRepo: categoryRepo,
	}
}

// Get product by id from repo
func (u *ProductUseCase) GetProductById(productId uint64) (*models.Product, error) {
	productById, err := u.ProductRepo.SelectProductById(productId)
	if err != nil {
		return nil, errors.ErrProductNotFound
	}

	categories, err := u.CategoryRepo.GetPathToCategory(productById.Category)
	if err != nil {
		return nil, errors.ErrCategoryNotFound
	}

	productById.CategoryPath = categories

	return productById, nil
}

// Get recommendations for product
func (u *ProductUseCase) GetProductRecommendationsById(productId uint64,
	paginator *models.PaginatorRecommendations) ([]*models.RecommendationProduct, error) {
	recommendationsByReviews, err := u.ProductRepo.SelectRecommendationsByReviews(productId, paginator.Count)
	if err != nil {
		return nil, errors.ErrProductNotFound
	}

	countMissingItems := paginator.Count - len(recommendationsByReviews)
	recommendationsByCategory := make([]*models.RecommendationProduct, 0)
	if countMissingItems > 0 {
		recommendationsByCategory, err = u.ProductRepo.SelectRecommendationsByCategory(productId, countMissingItems)
		if err != nil {
			return nil, errors.ErrProductNotFound
		}
	}

	return append(recommendationsByReviews, recommendationsByCategory...), nil
}

// Get range products by paginator settings from repo
func (u *ProductUseCase) GetRangeProducts(paginator *models.PaginatorProducts) (*models.RangeProducts, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	var filterString string
	if paginator.Filter != nil {
		filterString = u.ProductRepo.CreateFilterString(paginator.Filter)
	}

	// Max count pages in catalog
	countPages, err := u.ProductRepo.GetCountPages(paginator.Category, paginator.Count, filterString)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Keys for sort items in catalog
	sortString, err := u.ProductRepo.CreateSortString(paginator.SortKey, paginator.SortDirection)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get range of products
	products, err := u.ProductRepo.SelectRangeProducts(paginator, sortString, filterString)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	return &models.RangeProducts{
		ListPreviewProducts: products,
		MaxCountPages:       countPages,
	}, nil
}

// Find range products by search settings from repo
func (u *ProductUseCase) SearchRangeProducts(searchQuery *models.SearchQuery) (*models.RangeProducts, error) {
	if searchQuery.PageNum < 1 || searchQuery.Count < 1 {
		return nil, errors.ErrIncorrectSearchQuery
	}

	var filterString string
	if searchQuery.Filter != nil {
		filterString = u.ProductRepo.CreateFilterString(searchQuery.Filter)
	}

	// Max count pages for this search
	countPages, err := u.ProductRepo.GetCountSearchPages(searchQuery.Category,
		searchQuery.Count, searchQuery.QueryString, filterString)
	if err != nil {
		return nil, errors.ErrIncorrectSearchQuery
	}

	// Keys for sort items in result of search
	sortString, err := u.ProductRepo.CreateSortString(searchQuery.SortKey, searchQuery.SortDirection)
	if err != nil {
		return nil, errors.ErrIncorrectSearchQuery
	}

	// Get range of products
	products, err := u.ProductRepo.SearchRangeProducts(searchQuery, sortString, filterString)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	return &models.RangeProducts{
		ListPreviewProducts: products,
		MaxCountPages:       countPages,
	}, nil
}
