package product

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product Repository

type Repository interface {
	SelectProductById(productId uint64) (*models.Product, error)
	GetCountPages(category uint64, count int, filterString string) (int, error)
	GetCountSearchPages(category uint64, count int, searchString, filterString string) (int, error)
	CreateSortString(sortKey, sortDirection string) (string, error)
	SelectRangeProducts(paginator *models.PaginatorProducts,
		sortString, filterString string) ([]*models.ViewProduct, error)
	CreateFilterString(filter *models.ProductFilter) string
	SearchRangeProducts(searchQuery *models.SearchQuery,
		sortString, filterString string) ([]*models.ViewProduct, error)
	SelectRecommendationsByReviews(productId uint64, count int) ([]*models.RecommendationProduct, error)
	SelectRecommendationsByCategory(productId uint64, count int) ([]*models.RecommendationProduct, error)
}
