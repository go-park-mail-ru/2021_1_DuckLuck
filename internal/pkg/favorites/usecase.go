package favorites

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/favorites UseCase

type UseCase interface {
	AddProductToFavorites(productId, userId uint64) error
	DeleteProductFromFavorites(productId, userId uint64) error
	GetRangeFavorites(paginator *models.PaginatorFavorites,
		userId uint64) (*models.RangeFavorites, error)
	GetUserFavorites(userId uint64) (*models.UserFavorites, error)
}
