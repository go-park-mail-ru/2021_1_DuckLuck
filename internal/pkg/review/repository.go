package review

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product Repository

type Repository interface {
	SelectStatisticsByProductId(productId uint64) (*models.ReviewStatistics, error)
	CheckReview(userId uint64, productId uint64) bool
	AddReview(userId uint64, review *models.Review) (uint64, error)
	GetCountPages(productId uint64, countOrdersOnPage int) (int, error)
	CreateSortString(sortKey, sortDirection string) (string, error)
	SelectRangeReviews(productId uint64, sortString string,
		paginator *models.PaginatorReviews) ([]*models.ViewReview, error)
}
