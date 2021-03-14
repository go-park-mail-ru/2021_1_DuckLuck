package product

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type ProductUseCase struct {
	ProductRepo product.Repository
}

func NewUseCase(repo product.Repository) product.UseCase {
	return &ProductUseCase{
		ProductRepo: repo,
	}
}

func (lr *ProductUseCase) GetProductById(productId uint64) (*models.Product, error) {
	return lr.ProductRepo.SelectProductById(productId)
}

func (lr *ProductUseCase) SelectRangeProducts(paginator *models.PaginatorProducts) (*models.RangeProducts, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	return lr.ProductRepo.SelectRangeProducts(paginator)
}
