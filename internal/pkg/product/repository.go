package product

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
)

type Repository interface {
	GetById(productId uint64) (*models.Product, error)
	GetPaginateProducts(paginator *models.PaginatorProducts) (*models.RangeProducts, error)
}
