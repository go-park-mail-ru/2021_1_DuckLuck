package favorites

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

type UseCase interface {
	AddProductToFavorites(productId, userId uint64) error
	DeleteProductFromFavorites(productId, userId uint64) error
	GetRangeFavorites(paginator *models.PaginatorFavorites,
		userId uint64) (*models.RangeFavorites, error)
}
