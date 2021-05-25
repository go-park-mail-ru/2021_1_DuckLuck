package favorites

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/favorites Repository

type Repository interface {
	AddProductToFavorites(productId, userId uint64) error
	DeleteProductFromFavorites(productId, userId uint64) error
	GetCountPages(userId uint64, count int) (int, error)
	CreateSortString(sortKey, sortDirection string) (string, error)
	SelectRangeFavorites(paginator *models.PaginatorFavorites,
		sortString string, userId uint64) ([]*models.ViewFavorite, error)
	GetUserFavorites(userId uint64) (*models.UserFavorites, error)
}
