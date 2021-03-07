package repository

import (
	"sort"
	"sync"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/models"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product"
	server_errors "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/errors"
)

type LocalRepository struct {
	data map[uint64]*models.Product
	mu   *sync.RWMutex
}

func NewSessionLocalRepository() product.Repository {
	return &LocalRepository{
		data: map[uint64]*models.Product{
			1: &models.Product{
				Id:          1,
				Title:       "Test1",
				Cost:        100,
				Rating:      1,
				Description: "Good item",
				Category:    "Home",
				Image:       "/product/test.png",
			},
		},
		mu:   &sync.RWMutex{},
	}
}

func (lr *LocalRepository) GetById(productId uint64) (*models.Product, error) {
	lr.mu.RLock()
	productById, ok := lr.data[productId]
	lr.mu.RUnlock()

	if !ok {
		return nil, server_errors.ErrProductNotFound
	}

	return productById, nil
}

func (lr *LocalRepository) GetPaginateProducts(paginator *models.PaginatorProducts) (*models.RangeProducts, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, server_errors.ErrIncorrectPaginator
	}

	countPages :=  len(lr.data) / paginator.Count
	if len(lr.data) % paginator.Count > 0 {
		countPages++
	}

	if countPages < paginator.PageNum {
		return nil, server_errors.ErrProductsIsEmpty
	}

	products := make([]*models.Product, 7)
	lr.mu.RLock()
	for _, item := range lr.data {
		products = append(products, item)
	}
	lr.mu.RUnlock()

	compare := func(i, j int) bool { return false }
	switch paginator.SortKey {
	case models.ProductsCostSort:
		if paginator.SortDirection {
			compare =  func(i, j int) bool {
				return  products[i].Cost < products[j].Cost
			}
		} else {
			compare =  func(i, j int) bool {
				return  products[i].Cost > products[j].Cost
			}
		}
	case models.ProductsRatingSort:
		if paginator.SortDirection {
			compare =  func(i, j int) bool {
				return  products[i].Rating < products[j].Rating
			}
		} else {
			compare =  func(i, j int) bool {
				return  products[i].Rating > products[j].Rating
			}
		}
	}
	sort.SliceStable(products, compare)

	leftBorder := (paginator.PageNum - 1) * paginator.Count
	rightBorder := paginator.PageNum * paginator.Count
	if paginator.PageNum == countPages {
		rightBorder = len(lr.data)
	}

	return &models.RangeProducts{
		ArrayOfProducts: products[leftBorder : rightBorder],
		MaxCountPages: countPages,
	}, nil
}
