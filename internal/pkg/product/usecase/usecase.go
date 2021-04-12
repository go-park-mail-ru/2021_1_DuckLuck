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
	return u.ProductRepo.SelectProductById(productId)
}

// Get range products by paginator settings from repo
func (u *ProductUseCase) GetRangeProducts(paginator *models.PaginatorProducts) (*models.RangeProducts, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	categoriesId, err := u.CategoryRepo.GetAllSubCategoriesId(paginator.Category)
	if err != nil {
		return nil, err
	}

	return u.ProductRepo.SelectRangeProducts(paginator, categoriesId)
}
