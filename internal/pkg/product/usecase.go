package product

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product UseCase

type UseCase interface {
	GetProductById(productId uint64) (*models.Product, error)
	GetRangeProducts(paginator *models.PaginatorProducts) (*models.RangeProducts, error)
}
